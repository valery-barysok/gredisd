# Juno Inc., Test Task: In-memory cache

Simple implementation of Redis-like in-memory cache

Desired features:
- [x] Key-value storage with string, lists, dict support
- [x] Per-key TTL
- [x] Operations:
  - Get
  - Set
  - Update
  - Remove
  - Keys
- [x] Custom operations(Get i element on list, get value by key from dict, etc)
- [x] Golang API client
- [x] Telnet-like/HTTP-like API protocol

- [x] Provide some tests, API spec, deployment docs without full coverage, just a few cases and some examples of telnet/http calls to the server. 

Optional features:
- [ ] persistence to disk/db
- [ ] scaling(on server-side or on client-side, up to you)
- [x] auth
- [ ] performance tests
