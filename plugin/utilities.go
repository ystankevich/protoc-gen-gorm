package plugin

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	gorm "github.com/infobloxopen/protoc-gen-gorm/options"
)

// retrieves the GormMessageOptions from a message
func getMessageOptions(message *generator.Descriptor) *gorm.GormMessageOptions {
	if message.Options == nil {
		return nil
	}
	v, err := proto.GetExtension(message.Options, gorm.E_Opts)
	if err != nil {
		return nil
	}
	opts, ok := v.(*gorm.GormMessageOptions)
	if !ok {
		return nil
	}
	return opts
}

func getFieldOptions(field *descriptor.FieldDescriptorProto) *gorm.GormFieldOptions {
	if field.Options == nil {
		return nil
	}
	v, err := proto.GetExtension(field.Options, gorm.E_Field)
	if err != nil {
		return nil
	}
	opts, ok := v.(*gorm.GormFieldOptions)
	if !ok {
		return nil
	}
	return opts
}

func generateTransactionHandling(p *OrmPlugin) {
	p.P(`var err error`)
	p.P(`var txn *`, p.lftPkgName, `.Transaction`)
	p.P(`txn, ok := `, p.lftPkgName, `.FromContext(ctx)`)
	p.P(`if !ok {`)

	p.P(`defer func() {`)
	p.P(`var terr error`)
	p.P(`if err != nil {`)
	p.P(`terr = txn.Rollback()`)
	p.P(`} else {`)
	p.P(`if terr = txn.Commit(); terr != nil {`)
	p.P(`err = status.Error(codes.Internal, "+failed to commit transaction")`)
	p.P(`}`)
	p.P(`}`)

	p.P(`if terr == nil {`)
	p.P(`return`)
	p.P(`}`)

	p.P(`st := status.Convert(err)`)
	p.P(`st, serr := st.WithDetails(errdetails.New(codes.Internal, "gorm", terr.Error()))`)
	p.P(`// do not override error if failed to attach details`)
	p.P(`if serr == nil {`)
	p.P(`err = st.Err()`)
	p.P(`}`)
	p.P(`return`)
	p.P(`}()`)

	p.P(`txn = `, p.lftPkgName, `.NewTransaction(db)`)
	p.P(`}`)
	p.P(`db = txn.Begin()`)

}
