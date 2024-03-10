package listener

import (
	"context"
	"github.com/streadway/amqp"
	"github.com/wangliujing/foundation-framework/rabbitmq"
	"github.com/wangliujing/foundation-framework/util/ustring"
	"github.com/wangliujing/sensitive-word/internal/config"
	"github.com/wangliujing/sensitive-word/internal/logic"
	"github.com/wangliujing/sensitive-word/internal/svc"
)

type WordChangeListener struct {
	queueName      string
	consumer       string
	serviceContext *svc.ServiceContext
}

func NewWordChangeListener(conf config.Config, serviceContext *svc.ServiceContext) rabbitmq.Listener {
	uuid := ustring.GetUUID()
	return &WordChangeListener{
		queueName:      "WordChange-" + uuid,
		consumer:       uuid,
		serviceContext: serviceContext,
	}
}

func (m *WordChangeListener) CoroutineNum() int {
	return 1
}

func (m *WordChangeListener) ConsumerParam() rabbitmq.ConsumerParam {
	return rabbitmq.ConsumerParam{
		Consumer:  m.consumer,
		AutoAck:   true,
		Exclusive: true,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}
}

func (m *WordChangeListener) AckParam() rabbitmq.AckParam {
	return rabbitmq.AckParam{Multipart: false, Requeue: false}
}

func (m *WordChangeListener) QueueDeclareParam() rabbitmq.QueueDeclareParam {
	return rabbitmq.QueueDeclareParam{
		QueueName:    m.queueName,
		Durable:      false,
		AutoDelete:   true,
		Exclusive:    true,
		NoWait:       false,
		RoutingKey:   "amz_sensitive_center.sensitive_update",
		ExchangeName: "amz_sensitive_center",
	}
}

func (m *WordChangeListener) OnDelivery(delivery *amqp.Delivery) error {
	sensitiveWordLogic := logic.NewSensitiveWordLogic(context.Background(), m.serviceContext)
	sensitiveWordLogic.InitTrie()
	return nil
}

/*type Trigger struct {
	timer    *time.Timer
	waitTime time.Duration
}

func NewTrigger(fun func(), waitTime time.Duration) *Trigger {
	trigger := &Trigger{
		timer:    time.NewTimer(waitTime),
		waitTime: waitTime,
	}
	go func() {
		for {
			select {
			case <-trigger.timer.C:
				fun()
			}
		}
	}()
	return trigger
}

func (t *Trigger) Reset() {
	t.timer.Reset(t.waitTime)
}*/
