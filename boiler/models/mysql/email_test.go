// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package mysql

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testEmails(t *testing.T) {
	t.Parallel()

	query := Emails()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testEmailsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Emails().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testEmailsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Emails().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Emails().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testEmailsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := EmailSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Emails().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testEmailsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := EmailExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Email exists: %s", err)
	}
	if !e {
		t.Errorf("Expected EmailExists to return true, but got false.")
	}
}

func testEmailsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	emailFound, err := FindEmail(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if emailFound == nil {
		t.Error("want a record, got nil")
	}
}

func testEmailsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Emails().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testEmailsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Emails().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testEmailsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	emailOne := &Email{}
	emailTwo := &Email{}
	if err = randomize.Struct(seed, emailOne, emailDBTypes, false, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}
	if err = randomize.Struct(seed, emailTwo, emailDBTypes, false, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = emailOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = emailTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Emails().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testEmailsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	emailOne := &Email{}
	emailTwo := &Email{}
	if err = randomize.Struct(seed, emailOne, emailDBTypes, false, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}
	if err = randomize.Struct(seed, emailTwo, emailDBTypes, false, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = emailOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = emailTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Emails().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func emailBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *Email) error {
	*o = Email{}
	return nil
}

func emailAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *Email) error {
	*o = Email{}
	return nil
}

func emailAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *Email) error {
	*o = Email{}
	return nil
}

func emailBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Email) error {
	*o = Email{}
	return nil
}

func emailAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *Email) error {
	*o = Email{}
	return nil
}

func emailBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Email) error {
	*o = Email{}
	return nil
}

func emailAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *Email) error {
	*o = Email{}
	return nil
}

func emailBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Email) error {
	*o = Email{}
	return nil
}

func emailAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *Email) error {
	*o = Email{}
	return nil
}

func testEmailsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &Email{}
	o := &Email{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, emailDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Email object: %s", err)
	}

	AddEmailHook(boil.BeforeInsertHook, emailBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	emailBeforeInsertHooks = []EmailHook{}

	AddEmailHook(boil.AfterInsertHook, emailAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	emailAfterInsertHooks = []EmailHook{}

	AddEmailHook(boil.AfterSelectHook, emailAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	emailAfterSelectHooks = []EmailHook{}

	AddEmailHook(boil.BeforeUpdateHook, emailBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	emailBeforeUpdateHooks = []EmailHook{}

	AddEmailHook(boil.AfterUpdateHook, emailAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	emailAfterUpdateHooks = []EmailHook{}

	AddEmailHook(boil.BeforeDeleteHook, emailBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	emailBeforeDeleteHooks = []EmailHook{}

	AddEmailHook(boil.AfterDeleteHook, emailAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	emailAfterDeleteHooks = []EmailHook{}

	AddEmailHook(boil.BeforeUpsertHook, emailBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	emailBeforeUpsertHooks = []EmailHook{}

	AddEmailHook(boil.AfterUpsertHook, emailAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	emailAfterUpsertHooks = []EmailHook{}
}

func testEmailsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Emails().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testEmailsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(emailColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Emails().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testEmailsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testEmailsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := EmailSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testEmailsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Emails().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	emailDBTypes = map[string]string{`ID`: `int`, `Recipients`: `text`, `CC`: `text`, `BCC`: `text`, `Sender`: `varchar`, `Subject`: `varchar`, `MimeType`: `enum('html','text')`, `Body`: `longtext`, `Status`: `enum('created','sent','error')`, `Error`: `varchar`, `DTCreated`: `datetime`, `DTUpdated`: `datetime`, `DTSent`: `datetime`}
	_            = bytes.MinRead
)

func testEmailsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(emailPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(emailColumns) == len(emailPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Emails().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, emailDBTypes, true, emailPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testEmailsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(emailColumns) == len(emailPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Email{}
	if err = randomize.Struct(seed, o, emailDBTypes, true, emailColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Emails().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, emailDBTypes, true, emailPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(emailColumns, emailPrimaryKeyColumns) {
		fields = emailColumns
	} else {
		fields = strmangle.SetComplement(
			emailColumns,
			emailPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := EmailSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testEmailsUpsert(t *testing.T) {
	t.Parallel()

	if len(emailColumns) == len(emailPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}
	if len(mySQLEmailUniqueColumns) == 0 {
		t.Skip("Skipping table with no unique columns to conflict on")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Email{}
	if err = randomize.Struct(seed, &o, emailDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Email: %s", err)
	}

	count, err := Emails().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, emailDBTypes, false, emailPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Email struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Email: %s", err)
	}

	count, err = Emails().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}