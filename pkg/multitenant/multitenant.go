package multitenant

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
)

type Multitenant interface {
	New(tenant string) error
	NewFromOptions(tenant string, options *aruba.Options) error
	Add(tenant string, client aruba.Client)
	Get(tenant string) (aruba.Client, bool)
	MustGet(tenant string) aruba.Client
	GetOrNil(tenant string) aruba.Client
	CleanUp(from time.Duration)
}

type entry struct {
	client    aruba.Client
	lastUsage time.Time
}

type multitenant struct {
	clients  map[string]*entry
	template *aruba.Options

	lock sync.RWMutex
}

var _ Multitenant = (*multitenant)(nil)

func New() Multitenant {
	return &multitenant{
		clients: make(map[string]*entry),
	}
}

func NewWithTemplate(template *aruba.Options) Multitenant {
	return &multitenant{
		clients:  make(map[string]*entry),
		template: template,
	}
}

func (m *multitenant) New(tenant string) error {
	if m.template == nil {
		return errors.New("template is missing - use the `NewFromOptions` method")
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	template := m.template.DeepCopy()

	c, err := aruba.NewClient(template)
	if err != nil {
		return err
	}

	m.clients[tenant] = &entry{client: c, lastUsage: time.Now()}

	return nil
}

func (m *multitenant) NewFromOptions(tenant string, options *aruba.Options) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	c, err := aruba.NewClient(options)
	if err != nil {
		return err
	}

	m.clients[tenant] = &entry{client: c, lastUsage: time.Now()}

	return nil
}

func (m *multitenant) Add(tenant string, client aruba.Client) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.clients[tenant] = &entry{client: client, lastUsage: time.Now()}
}

func (m *multitenant) Get(tenant string) (aruba.Client, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	e, ok := m.clients[tenant]

	if ok {
		e.lastUsage = time.Now()
	}

	return e.client, ok
}

func (m *multitenant) MustGet(tenant string) aruba.Client {
	m.lock.RLock()
	defer m.lock.RUnlock()

	e, ok := m.clients[tenant]
	if !ok {
		log.Fatalf("client for tenant '%s' not found", tenant)
	}

	e.lastUsage = time.Now()

	return e.client
}

func (m *multitenant) GetOrNil(tenant string) aruba.Client {
	m.lock.RLock()
	defer m.lock.RUnlock()

	e, ok := m.clients[tenant]
	if !ok {
		return nil
	}

	e.lastUsage = time.Now()

	return e.client
}

func (m *multitenant) CleanUp(from time.Duration) {
	m.lock.Lock()
	defer m.lock.Unlock()

	cleanupTime := time.Now().Add(-1 * from)

	for t, e := range m.clients {
		if e.lastUsage.Before(cleanupTime) {
			delete(m.clients, t)
		}
	}
}
