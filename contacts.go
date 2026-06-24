// Package contacts is the togo contacts base plugin. It defines a normalized
// Contact model + a ContactsProvider driver interface; provider plugins (e.g.
// contacts-google) call contacts.RegisterDriver and depend on this package.
//
// Blank-import a provider and set CONTACTS_DRIVER to select it. The default
// "null" driver returns no contacts (safe for dev/tests).
package contacts

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/togo-framework/togo"
)

// Contact is a normalized contact record (provider-agnostic).
type Contact struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Emails []string          `json:"emails,omitempty"`
	Phones []string          `json:"phones,omitempty"`
	Photo  string            `json:"photo,omitempty"`
	Org    string            `json:"org,omitempty"`
	Source string            `json:"source,omitempty"`
	Extra  map[string]string `json:"extra,omitempty"`
}

// ContactsProvider imports contacts from an external source (Google, a CRM, …).
type ContactsProvider interface {
	// List returns a page of contacts. pageToken "" starts at the beginning; the
	// returned next token is "" when there are no more pages.
	List(ctx context.Context, pageToken string) (contacts []Contact, next string, err error)
	// Get returns one contact by its provider id.
	Get(ctx context.Context, id string) (*Contact, error)
}

// DriverFactory builds a provider from the kernel (env/config).
type DriverFactory func(k *togo.Kernel) (ContactsProvider, error)

var (
	regMu   sync.RWMutex
	drivers = map[string]DriverFactory{}
)

// RegisterDriver registers a contacts driver by name (call from a plugin's init()).
func RegisterDriver(name string, f DriverFactory) {
	regMu.Lock()
	drivers[name] = f
	regMu.Unlock()
}

func init() {
	RegisterDriver("null", func(k *togo.Kernel) (ContactsProvider, error) { return nullProvider{}, nil })

	togo.RegisterProviderFunc("contacts", togo.PriorityService, func(k *togo.Kernel) error {
		name := os.Getenv("CONTACTS_DRIVER")
		if name == "" {
			name = "null"
		}
		regMu.RLock()
		f, ok := drivers[name]
		regMu.RUnlock()
		if !ok {
			return fmt.Errorf("contacts: unknown driver %q (install its plugin, e.g. togo install togo-framework/contacts-%s)", name, name)
		}
		p, err := f(k)
		if err != nil {
			return err
		}
		k.Set("contacts", &Service{provider: p, driver: name})
		return nil
	})
}

// Service is the contacts runtime stored on the kernel (k.Get("contacts")).
type Service struct {
	provider ContactsProvider
	driver   string
}

// Driver returns the active driver name.
func (s *Service) Driver() string { return s.driver }

// Provider returns the active driver implementation.
func (s *Service) Provider() ContactsProvider { return s.provider }

// Get returns one contact by provider id.
func (s *Service) Get(ctx context.Context, id string) (*Contact, error) { return s.provider.Get(ctx, id) }

// Sync pulls every contact from the provider, paging through all pages.
func (s *Service) Sync(ctx context.Context) ([]Contact, error) {
	var all []Contact
	token := ""
	for {
		page, next, err := s.provider.List(ctx, token)
		if err != nil {
			return all, err
		}
		all = append(all, page...)
		if next == "" {
			break
		}
		token = next
	}
	return all, nil
}

// FromKernel fetches the contacts service from the kernel container.
func FromKernel(k *togo.Kernel) (*Service, bool) {
	v, ok := k.Get("contacts")
	if !ok {
		return nil, false
	}
	s, ok := v.(*Service)
	return s, ok
}

// nullProvider is the safe default: no contacts.
type nullProvider struct{}

func (nullProvider) List(context.Context, string) ([]Contact, string, error) { return nil, "", nil }
func (nullProvider) Get(context.Context, string) (*Contact, error) {
	return nil, fmt.Errorf("contacts: no driver configured (set CONTACTS_DRIVER + install a provider like togo-framework/contacts-google)")
}
