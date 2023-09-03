package audited

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type auditableInterface interface {
	SetCreatedBy(createdBy interface{})
	GetCreatedBy() string
	SetUpdatedBy(updatedBy interface{})
	GetUpdatedBy() string
}

func isAuditable(scope *gorm.DB) (isAuditable bool) {
	// TODO: check if scope.Statement.Schema.ModelType is able to be null as it is a reflect.Type
	if scope.Statement.Schema.ModelType == nil {
		return false
	}

	// TODO: check if scope.Statement.Schema.ModelType is correcthere
	_, isAuditable = reflect.New(scope.Statement.Schema.ModelType).Interface().(auditableInterface)
	return
}

func getCurrentUser(scope *gorm.DB) (string, bool) {
	var user interface{}
	var hasUser bool

	user, hasUser = scope.Get("audited:current_user")

	if !hasUser {
		user, hasUser = scope.Get("qor:current_user")
	}

	if hasUser {
		var currentUser string
		// TODO: identify what the new function to get the primary field of a type is
		// PrioritizedPrimaryField or PrimaryFields candidates may need to define further as there are composite primary key options here
		// Many need to indirect and get the field value of the primary field where primary key is not already captured
		/*
			var currentUser string
			if primaryField := scope.New(user).PrimaryField(); primaryField != nil {
				currentUser = fmt.Sprintf("%v", primaryField.Field.Interface())
			} else {
				currentUser = fmt.Sprintf("%v", user)
			}
		*/

		// Create a new instance of the type of the model user
		db2 := scope.Session(&gorm.Session{NewDB: true}).Model(user)
		if primaryField := db2.Statement.Schema.PrioritizedPrimaryField; primaryField != nil {
			currentUser = fmt.Sprintf("%v", reflect.Indirect(reflect.ValueOf(user)).FieldByName(primaryField.Name).Interface())
		} else {
			currentUser = fmt.Sprintf("%v", user)
		}

		return currentUser, true
	}

	return "", false
}

func assignCreatedBy(scope *gorm.DB) {
	if isAuditable(scope) {
		if user, ok := getCurrentUser(scope); ok {
			scope.Statement.SetColumn("CreatedBy", user)
		}
	}
}

func assignUpdatedBy(scope *gorm.DB) {
	if isAuditable(scope) {
		if user, ok := getCurrentUser(scope); ok {
			if attrs, ok := scope.InstanceGet("gorm:update_attrs"); ok {
				updateAttrs := attrs.(map[string]interface{})
				updateAttrs["updated_by"] = user
				scope.InstanceSet("gorm:update_attrs", updateAttrs)
			} else {
				scope.Statement.SetColumn("UpdatedBy", user)
			}
		}
	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()
	if callback.Create().Get("audited:assign_created_by") == nil {
		callback.Create().After("gorm:before_create").Register("audited:assign_created_by", assignCreatedBy)
	}
	if callback.Update().Get("audited:assign_updated_by") == nil {
		callback.Update().After("gorm:before_update").Register("audited:assign_updated_by", assignUpdatedBy)
	}
}
