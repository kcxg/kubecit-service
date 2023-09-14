package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"kubecit-service/ent"
	"kubecit-service/ent/course"
	"kubecit-service/internal/biz"
	"math/rand"
	"time"
)

type orderRepo struct {
	data *Data
	log  *log.Helper
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (or *orderRepo) Create(ctx context.Context, courseIds []int32) (*biz.Order, error) {
	result, err := or.data.WithResultTx(ctx, func(tx *ent.Tx) (interface{}, error) {
		orderTx, err := or.createOrderTx(ctx, courseIds)
		if err != nil {
			return nil, err
		}
		return orderTx, nil
	})
	if err != nil {
		return nil, err
	}
	if order, ok := result.(*biz.Order); ok {
		return order, nil
	}
	return nil, err
}

func (or *orderRepo) createOrderTx(ctx context.Context, courseIds []int32) (*biz.Order, error) {
	var coursePrice int32
	for _, id := range courseIds {
		c, err := or.data.db.Course.Query().Where(course.IDEQ(int(id))).First(ctx)
		if err != nil {
			return nil, errors.New(400, "课程不存在！", "根据课程ID查询课程失败")
		}
		coursePrice += int32(c.Price)
	}
	userIdRaw := ctx.Value("userId")
	userId, ok := userIdRaw.(int32)
	if !ok {
		return nil, errors.New(400, "用户ID不存在", "从token中解析不出用户ID")
	}
	// Proceed with the userId value
	orderObj, err := or.data.db.Orders.Create().SetOrderSn(GenerateOrderSn(userId)).SetTradePrice(coursePrice).SetUserID(userId).Save(ctx)
	if err != nil {
		return nil, errors.BadRequest("创建订单失败", err.Error())
	}

	infos := make([]*biz.OrderInfo, 0)
	for _, id := range courseIds {
		c, err := or.data.db.Course.Query().Where(course.IDEQ(int(id))).First(ctx)
		if err != nil {
			return nil, errors.BadRequest("根据课程ID查询课程失败", err.Error())
		}
		info, err := or.data.db.OrderInfos.Create().SetOrderID(int32(orderObj.ID)).SetCourseID(id).SetCourseDescribe(c.Detail).SetCourseName(c.Name).SetCoursePrice(int32(c.Price)).Save(ctx)
		if err != nil {
			return nil, errors.BadRequest("创建订单详情失败", err.Error())
		}

		infos = append(infos, &biz.OrderInfo{
			Id:             info.ID,
			CourseName:     info.CourseName,
			CourseDescribe: info.CourseDescribe,
			CoursePrice:    int32(c.Price),
			CreateTime:     info.CreateTime,
			UpdateTime:     info.UpdateTime,
			OrderId:        int32(orderObj.ID),
			CourseId:       id,
		})
	}
	order := &biz.Order{
		Id:         orderObj.ID,
		UserId:     userId,
		OrderSn:    orderObj.OrderSn,
		PayType:    orderObj.PayType,
		PayStatus:  orderObj.PayStatus,
		TradePrice: orderObj.TradePrice,
		TradeNo:    orderObj.TradeNo,
		PayTime:    orderObj.PayTime,
		CreateTime: orderObj.CreateTime,
		UpdateTime: orderObj.UpdateTime,
		Info:       infos,
	}
	return order, nil
}

func GenerateOrderSn(userId int32) string {
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(), userId, rand.Intn(90)+10)
	return orderSn
}
