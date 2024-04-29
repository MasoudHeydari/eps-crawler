// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/karust/openserp/ent/serp"
)

// SERPCreate is the builder for creating a SERP entity.
type SERPCreate struct {
	config
	mutation *SERPMutation
	hooks    []Hook
}

// SetURL sets the "url" field.
func (sc *SERPCreate) SetURL(s string) *SERPCreate {
	sc.mutation.SetURL(s)
	return sc
}

// SetTitle sets the "title" field.
func (sc *SERPCreate) SetTitle(s string) *SERPCreate {
	sc.mutation.SetTitle(s)
	return sc
}

// SetDescription sets the "description" field.
func (sc *SERPCreate) SetDescription(s string) *SERPCreate {
	sc.mutation.SetDescription(s)
	return sc
}

// SetLocation sets the "location" field.
func (sc *SERPCreate) SetLocation(s string) *SERPCreate {
	sc.mutation.SetLocation(s)
	return sc
}

// SetContactInfo sets the "contact_info" field.
func (sc *SERPCreate) SetContactInfo(s []string) *SERPCreate {
	sc.mutation.SetContactInfo(s)
	return sc
}

// SetKeyWords sets the "key_words" field.
func (sc *SERPCreate) SetKeyWords(s []string) *SERPCreate {
	sc.mutation.SetKeyWords(s)
	return sc
}

// SetIsRead sets the "is_read" field.
func (sc *SERPCreate) SetIsRead(b bool) *SERPCreate {
	sc.mutation.SetIsRead(b)
	return sc
}

// SetNillableIsRead sets the "is_read" field if the given value is not nil.
func (sc *SERPCreate) SetNillableIsRead(b *bool) *SERPCreate {
	if b != nil {
		sc.SetIsRead(*b)
	}
	return sc
}

// SetCreatedAt sets the "created_at" field.
func (sc *SERPCreate) SetCreatedAt(t time.Time) *SERPCreate {
	sc.mutation.SetCreatedAt(t)
	return sc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (sc *SERPCreate) SetNillableCreatedAt(t *time.Time) *SERPCreate {
	if t != nil {
		sc.SetCreatedAt(*t)
	}
	return sc
}

// Mutation returns the SERPMutation object of the builder.
func (sc *SERPCreate) Mutation() *SERPMutation {
	return sc.mutation
}

// Save creates the SERP in the database.
func (sc *SERPCreate) Save(ctx context.Context) (*SERP, error) {
	sc.defaults()
	return withHooks(ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SERPCreate) SaveX(ctx context.Context) *SERP {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SERPCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SERPCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sc *SERPCreate) defaults() {
	if _, ok := sc.mutation.IsRead(); !ok {
		v := serp.DefaultIsRead
		sc.mutation.SetIsRead(v)
	}
	if _, ok := sc.mutation.CreatedAt(); !ok {
		v := serp.DefaultCreatedAt()
		sc.mutation.SetCreatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SERPCreate) check() error {
	if _, ok := sc.mutation.URL(); !ok {
		return &ValidationError{Name: "url", err: errors.New(`ent: missing required field "SERP.url"`)}
	}
	if v, ok := sc.mutation.URL(); ok {
		if err := serp.URLValidator(v); err != nil {
			return &ValidationError{Name: "url", err: fmt.Errorf(`ent: validator failed for field "SERP.url": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New(`ent: missing required field "SERP.title"`)}
	}
	if v, ok := sc.mutation.Title(); ok {
		if err := serp.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "SERP.title": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New(`ent: missing required field "SERP.description"`)}
	}
	if _, ok := sc.mutation.Location(); !ok {
		return &ValidationError{Name: "location", err: errors.New(`ent: missing required field "SERP.location"`)}
	}
	if _, ok := sc.mutation.ContactInfo(); !ok {
		return &ValidationError{Name: "contact_info", err: errors.New(`ent: missing required field "SERP.contact_info"`)}
	}
	if _, ok := sc.mutation.KeyWords(); !ok {
		return &ValidationError{Name: "key_words", err: errors.New(`ent: missing required field "SERP.key_words"`)}
	}
	if _, ok := sc.mutation.IsRead(); !ok {
		return &ValidationError{Name: "is_read", err: errors.New(`ent: missing required field "SERP.is_read"`)}
	}
	if _, ok := sc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "SERP.created_at"`)}
	}
	return nil
}

func (sc *SERPCreate) sqlSave(ctx context.Context) (*SERP, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *SERPCreate) createSpec() (*SERP, *sqlgraph.CreateSpec) {
	var (
		_node = &SERP{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(serp.Table, sqlgraph.NewFieldSpec(serp.FieldID, field.TypeInt))
	)
	if value, ok := sc.mutation.URL(); ok {
		_spec.SetField(serp.FieldURL, field.TypeString, value)
		_node.URL = value
	}
	if value, ok := sc.mutation.Title(); ok {
		_spec.SetField(serp.FieldTitle, field.TypeString, value)
		_node.Title = value
	}
	if value, ok := sc.mutation.Description(); ok {
		_spec.SetField(serp.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := sc.mutation.Location(); ok {
		_spec.SetField(serp.FieldLocation, field.TypeString, value)
		_node.Location = value
	}
	if value, ok := sc.mutation.ContactInfo(); ok {
		_spec.SetField(serp.FieldContactInfo, field.TypeJSON, value)
		_node.ContactInfo = value
	}
	if value, ok := sc.mutation.KeyWords(); ok {
		_spec.SetField(serp.FieldKeyWords, field.TypeJSON, value)
		_node.KeyWords = value
	}
	if value, ok := sc.mutation.IsRead(); ok {
		_spec.SetField(serp.FieldIsRead, field.TypeBool, value)
		_node.IsRead = value
	}
	if value, ok := sc.mutation.CreatedAt(); ok {
		_spec.SetField(serp.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	return _node, _spec
}

// SERPCreateBulk is the builder for creating many SERP entities in bulk.
type SERPCreateBulk struct {
	config
	err      error
	builders []*SERPCreate
}

// Save creates the SERP entities in the database.
func (scb *SERPCreateBulk) Save(ctx context.Context) ([]*SERP, error) {
	if scb.err != nil {
		return nil, scb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*SERP, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SERPMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SERPCreateBulk) SaveX(ctx context.Context) []*SERP {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SERPCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SERPCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}
