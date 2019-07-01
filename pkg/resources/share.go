package resources

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var shareProperties = []string{
	"comment",
}

var shareSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the share; must be unique for the account in which the share is created.",
	},
	"comment": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the managed account.",
	},
	"accounts": &schema.Schema{
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "A list of accounts to be added to the share.",
	},
}

// Share returns a pointer to the resource representing a share
func Share() *schema.Resource {
	return &schema.Resource{
		Create: CreateShare,
		Read:   ReadShare,
		Update: UpdateShare,
		Delete: DeleteShare,
		Exists: ShareExists,

		Schema: shareSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateShare implements schema.CreateFunc
func CreateShare(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	builder := snowflake.Share(name).Create()
	builder.SetString("COMMENT", data.Get("comment").(string))

	err := DBExec(db, builder.Statement())
	if err != nil {
		return errors.Wrapf(err, "error creating share")
	}
	data.SetId(name)

	// Adding accounts must be done via an ALTER query

	// @TODO flesh out the share type in the snowflake package since it doesn't
	// follow the normal generic rules
	err = setAccounts(data, meta)
	if err != nil {
		return err
	}

	return ReadShare(data, meta)
}

func setAccounts(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)
	accs := expandStringList(data.Get("accounts").(*schema.Set).List())

	if len(accs) > 0 {
		q := fmt.Sprintf(`ALTER SHARE "%v" SET ACCOUNTS=%v`, name, strings.Join(accs, ","))
		err := DBExec(db, q)
		if err != nil {
			return errors.Wrapf(err, "error adding accounts to share %v", name)
		}
	}

	return nil
}

// ReadShare implements schema.ReadFunc
func ReadShare(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.Share(id).Show()
	row := db.QueryRow(stmt)

	var createdOn, kind, name, databaseName, to, owner, comment sql.NullString
	err := row.Scan(&createdOn, &kind, &name, &databaseName, &to, &owner, &comment)
	if err != nil {
		return err
	}

	// TODO turn this into a loop after we switch to scanning in a struct
	err = data.Set("name", StripAccountFromName(name.String))
	if err != nil {
		return err
	}
	err = data.Set("comment", comment.String)
	if err != nil {
		return err
	}

	accs := strings.Split(to.String, ", ")
	err = data.Set("accounts", accs)

	return err
}

// UpdateShare implements schema.UpdateFunc
func UpdateShare(data *schema.ResourceData, meta interface{}) error {
	// Change the accounts first - this is a special case and won't work using the generic method
	if data.HasChange("accounts") {
		err := setAccounts(data, meta)
		if err != nil {
			return err
		}
	}

	return UpdateResource("this does not seem to be used", shareProperties, shareSchema, snowflake.Share, ReadShare)(data, meta)
}

// DeleteShare implements schema.DeleteFunc
func DeleteShare(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("this does not seem to be used", snowflake.Share)(data, meta)
}

// ShareExists implements schema.ExistsFunc
func ShareExists(data *schema.ResourceData, meta interface{}) (bool, error) {
	db := meta.(*sql.DB)
	id := data.Id()

	stmt := snowflake.Share(id).Show()
	rows, err := db.Query(stmt)
	if err != nil {
		return false, err
	}

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

// StripAccountFromName removes the accout prefix from a resource (e.g. a share)
// that returns it (e.g. yt12345.my_share should just be my_share)
func StripAccountFromName(s string) string {
	return s[strings.Index(s, ".")+1:]
}
