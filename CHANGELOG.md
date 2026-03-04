# Changelog

## 1.0.0 (2026-03-04)


### Features

* add automatic versioning and release infrastructure ([8a0a3b8](https://github.com/danjdewhurst/envio/commit/8a0a3b87b7e6d63f0c710e7421aa73fafaf366dd))
* add HTTPS support with mkcert for local .test domains ([f1f7994](https://github.com/danjdewhurst/envio/commit/f1f79945db58479c59d3b6f7c426195553f27e71))
* add local .test domain support via shared Traefik reverse proxy ([a7e4946](https://github.com/danjdewhurst/envio/commit/a7e49464d585ada5bc334203699aad1874e9a7e5))


### Bug Fixes

* add Dockerfile build for FrankenPHP variant with PHP extensions ([a8e2f46](https://github.com/danjdewhurst/envio/commit/a8e2f46af2cbc3a49b2c26bf4fb2bd532a981ca9))
* use literal DB credentials to prevent .env interpolation conflicts ([5c9f701](https://github.com/danjdewhurst/envio/commit/5c9f701ea20ceff3aac9514513fda052a81eafab))


### Refactoring

* remove unused setup-dns command ([b6a48dc](https://github.com/danjdewhurst/envio/commit/b6a48dcc629d9c5dd303060352a64085d7e4f25d))
