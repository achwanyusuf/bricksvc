// Code generated by SQLBoiler 4.16.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package entity

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testTransferJobs(t *testing.T) {
	t.Parallel()

	query := TransferJobs()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testTransferJobsSoftDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx, false); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTransferJobsQuerySoftDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := TransferJobs().DeleteAll(ctx, tx, false); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTransferJobsSliceSoftDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TransferJobSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx, false); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTransferJobsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx, true); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTransferJobsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := TransferJobs().DeleteAll(ctx, tx, true); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTransferJobsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TransferJobSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx, true); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testTransferJobsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := TransferJobExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if TransferJob exists: %s", err)
	}
	if !e {
		t.Errorf("Expected TransferJobExists to return true, but got false.")
	}
}

func testTransferJobsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	transferJobFound, err := FindTransferJob(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if transferJobFound == nil {
		t.Error("want a record, got nil")
	}
}

func testTransferJobsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = TransferJobs().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testTransferJobsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := TransferJobs().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testTransferJobsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	transferJobOne := &TransferJob{}
	transferJobTwo := &TransferJob{}
	if err = randomize.Struct(seed, transferJobOne, transferJobDBTypes, false, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}
	if err = randomize.Struct(seed, transferJobTwo, transferJobDBTypes, false, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = transferJobOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = transferJobTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := TransferJobs().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testTransferJobsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	transferJobOne := &TransferJob{}
	transferJobTwo := &TransferJob{}
	if err = randomize.Struct(seed, transferJobOne, transferJobDBTypes, false, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}
	if err = randomize.Struct(seed, transferJobTwo, transferJobDBTypes, false, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = transferJobOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = transferJobTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func transferJobBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *TransferJob) error {
	*o = TransferJob{}
	return nil
}

func transferJobAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *TransferJob) error {
	*o = TransferJob{}
	return nil
}

func transferJobAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *TransferJob) error {
	*o = TransferJob{}
	return nil
}

func transferJobBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *TransferJob) error {
	*o = TransferJob{}
	return nil
}

func transferJobAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *TransferJob) error {
	*o = TransferJob{}
	return nil
}

func transferJobBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *TransferJob) error {
	*o = TransferJob{}
	return nil
}

func transferJobAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *TransferJob) error {
	*o = TransferJob{}
	return nil
}

func transferJobBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *TransferJob) error {
	*o = TransferJob{}
	return nil
}

func transferJobAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *TransferJob) error {
	*o = TransferJob{}
	return nil
}

func testTransferJobsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &TransferJob{}
	o := &TransferJob{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, transferJobDBTypes, false); err != nil {
		t.Errorf("Unable to randomize TransferJob object: %s", err)
	}

	AddTransferJobHook(boil.BeforeInsertHook, transferJobBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	transferJobBeforeInsertHooks = []TransferJobHook{}

	AddTransferJobHook(boil.AfterInsertHook, transferJobAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	transferJobAfterInsertHooks = []TransferJobHook{}

	AddTransferJobHook(boil.AfterSelectHook, transferJobAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	transferJobAfterSelectHooks = []TransferJobHook{}

	AddTransferJobHook(boil.BeforeUpdateHook, transferJobBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	transferJobBeforeUpdateHooks = []TransferJobHook{}

	AddTransferJobHook(boil.AfterUpdateHook, transferJobAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	transferJobAfterUpdateHooks = []TransferJobHook{}

	AddTransferJobHook(boil.BeforeDeleteHook, transferJobBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	transferJobBeforeDeleteHooks = []TransferJobHook{}

	AddTransferJobHook(boil.AfterDeleteHook, transferJobAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	transferJobAfterDeleteHooks = []TransferJobHook{}

	AddTransferJobHook(boil.BeforeUpsertHook, transferJobBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	transferJobBeforeUpsertHooks = []TransferJobHook{}

	AddTransferJobHook(boil.AfterUpsertHook, transferJobAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	transferJobAfterUpsertHooks = []TransferJobHook{}
}

func testTransferJobsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTransferJobsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(transferJobColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testTransferJobsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
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

func testTransferJobsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := TransferJobSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testTransferJobsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := TransferJobs().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	transferJobDBTypes = map[string]string{`ID`: `integer`, `JobID`: `character varying`, `APIKey`: `text`, `Payload`: `jsonb`, `Status`: `enum.transferstatus('failed','pending','success')`, `CreatedBy`: `integer`, `CreatedAt`: `timestamp with time zone`, `UpdatedBy`: `integer`, `UpdatedAt`: `timestamp with time zone`, `DeletedBy`: `integer`, `DeletedAt`: `timestamp with time zone`}
	_                  = bytes.MinRead
)

func testTransferJobsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(transferJobPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(transferJobAllColumns) == len(transferJobPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testTransferJobsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(transferJobAllColumns) == len(transferJobPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &TransferJob{}
	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, transferJobDBTypes, true, transferJobPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(transferJobAllColumns, transferJobPrimaryKeyColumns) {
		fields = transferJobAllColumns
	} else {
		fields = strmangle.SetComplement(
			transferJobAllColumns,
			transferJobPrimaryKeyColumns,
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

	slice := TransferJobSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testTransferJobsUpsert(t *testing.T) {
	t.Parallel()

	if len(transferJobAllColumns) == len(transferJobPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := TransferJob{}
	if err = randomize.Struct(seed, &o, transferJobDBTypes, true); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert TransferJob: %s", err)
	}

	count, err := TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, transferJobDBTypes, false, transferJobPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize TransferJob struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert TransferJob: %s", err)
	}

	count, err = TransferJobs().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
