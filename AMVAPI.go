package arn

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ Publishable            = (*AMV)(nil)
	_ Likeable               = (*AMV)(nil)
	_ LikeEventReceiver      = (*AMV)(nil)
	_ fmt.Stringer           = (*AMV)(nil)
	_ api.Newable            = (*AMV)(nil)
	_ api.Editable           = (*AMV)(nil)
	_ api.Deletable          = (*AMV)(nil)
	_ api.ArrayEventListener = (*AMV)(nil)
)

// Actions
func init() {
	API.RegisterActions("AMV", []*api.Action{
		// Publish
		PublishAction(),

		// Unpublish
		UnpublishAction(),

		// Like
		LikeAction(),

		// Unlike
		UnlikeAction(),
	})
}

// Create sets the data for a new AMV with data we received from the API request.
func (amv *AMV) Create(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	amv.ID = GenerateID("AMV")
	amv.Created = DateTimeUTC()
	amv.CreatedBy = user.ID

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "create", "AMV", amv.ID, "", "", "")
	logEntry.Save()

	return amv.Unpublish()
}

// Edit updates the external media object.
func (amv *AMV) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "AMV", amv.ID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	return false, nil
}

// OnAppend saves a log entry.
func (amv *AMV) OnAppend(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayAppend", "AMV", amv.ID, fmt.Sprintf("%s[%d]", key, index), "", fmt.Sprint(obj))
	logEntry.Save()
}

// OnRemove saves a log entry.
func (amv *AMV) OnRemove(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayRemove", "AMV", amv.ID, fmt.Sprintf("%s[%d]", key, index), fmt.Sprint(obj), "")
	logEntry.Save()
}

// AfterEdit updates the metadata.
func (amv *AMV) AfterEdit(ctx *aero.Context) error {
	amv.Edited = DateTimeUTC()
	amv.EditedBy = GetUserFromContext(ctx).ID
	return nil
}

// DeleteInContext deletes the amv in the given context.
func (amv *AMV) DeleteInContext(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "delete", "AMV", amv.ID, "", fmt.Sprint(amv), "")
	logEntry.Save()

	return amv.Delete()
}

// Delete deletes the object from the database.
func (amv *AMV) Delete() error {
	if amv.IsDraft {
		draftIndex := amv.Creator().DraftIndex()
		draftIndex.AMVID = ""
		draftIndex.Save()
	}

	if amv.File != "" {
		err := os.Remove(path.Join(Root, "videos", "amvs", amv.File))

		if err != nil {
			return err
		}
	}

	DB.Delete("AMV", amv.ID)
	return nil
}

// Authorize returns an error if the given API POST request is not authorized.
func (amv *AMV) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	if action == "delete" {
		if user.Role != "editor" && user.Role != "admin" {
			return errors.New("Insufficient permissions")
		}
	}

	return nil
}

// Save saves the amv object in the database.
func (amv *AMV) Save() {
	DB.Set("AMV", amv.ID, amv)
}
