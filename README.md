# contacts — togo contacts plugin

The **contacts** base plugin for [togo](https://to-go.dev). It defines a normalized
`Contact` model and a `ContactsProvider` driver interface so you can import/sync
contacts from anywhere — Google, a CRM, … — behind one API.

```bash
togo install togo-framework/contacts
togo install togo-framework/contacts-google   # a provider
```

Select the provider with `CONTACTS_DRIVER` (default `null` — no contacts):

```env
CONTACTS_DRIVER=google
```

## Use it

```go
import "github.com/togo-framework/contacts"

svc, _ := contacts.FromKernel(k)
all, err := svc.Sync(ctx)        // pull every contact from the provider
c, err  := svc.Get(ctx, id)      // one contact by provider id
```

## Write a provider

A provider is a tiny module that registers a driver in `init()`:

```go
func init() {
    contacts.RegisterDriver("mycrm", func(k *togo.Kernel) (contacts.ContactsProvider, error) {
        return &provider{/* read env */}, nil
    })
}
```

It implements `List(ctx, pageToken) ([]Contact, next, err)` and `Get(ctx, id)`.
`contacts-google` is the reference provider (Google People API); the same shape
works for any CRM.

MIT
