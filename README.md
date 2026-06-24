<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/contacts</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/contacts"><img src="https://pkg.go.dev/badge/github.com/togo-framework/contacts.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/contacts
```

<!-- /togo-header -->

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

<!-- togo-sponsors -->
---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
