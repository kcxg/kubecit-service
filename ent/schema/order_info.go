package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"time"
)

// OrderInfos holds the schema definition for the OrderInfos entity.
type OrderInfos struct {
	ent.Schema
}

// Fields of the OrderInfos.
func (OrderInfos) Fields() []ent.Field {

	return []ent.Field{

		field.Int32("order_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("	订单id"),

		field.Int32("course_id").SchemaType(map[string]string{
			dialect.MySQL: "int", // Override MySQL.
		}).Comment("课程id"),

		field.String("course_name").SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(64)", // Override MySQL.
		}).Comment("课程名称"),

		field.Float("course_price").SchemaType(map[string]string{
			dialect.MySQL: "decimal", // Override MySQL.
		}).Comment("课程价格"),

		field.Text("course_describe").SchemaType(map[string]string{
			dialect.MySQL: "text", // Override MySQL.
		}).Comment("课程描述"),

		field.Time("create_time").SchemaType(map[string]string{
			dialect.MySQL: "datetime", // Override MySQL.
		}).Default(time.Now()).Comment("创建时间"),

		field.Time("update_time").SchemaType(map[string]string{
			dialect.MySQL: "datetime", // Override MySQL.
		}).Default(time.Now()).UpdateDefault(time.Now).Comment("更新时间"),
	}
}

// Edges of the OrderInfos.
func (OrderInfos) Edges() []ent.Edge {
	return nil
}