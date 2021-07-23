# 👷 CloudFlare Workers in Go <img src="https://github.com/fikisipi/cloudflare-workers-go/actions/workflows/main.yml/badge.svg" />

`cfgo` uses WebAssembly to run Go projects as CF Workers. It exposes the APIs,
patches the missing runtime functions and glues the compiler to the CloudFlare tools.

To set up a project, install [CloudFlare Wrangler](https://github.com/cloudflare/wrangler) and run:

```
wrangler generate yourapp https://github.com/fikisipi/cloudflare-workers-go
```
### 🚴 Example and deployment
A demo with request handling is available in  `src/main.go`.

Run it using `wrangler dev`. To deploy
live, use `wrangler publish`.

### 🚧️ TODO
* [x] Event/Request handling API
* [x] fetch API
* [x] Handle wasm_exec from non-latest (<1.16) & tinygo 
* [x] KV for Workers API
   * TODO : add metadata and cursor pagination
* [ ] WebSocket API
* [ ] Support for streaming & bytes in fetch
* [ ] 💥 reducing worker size
   * code stripping? (already doing AST optimization in `wasm_exec`)
   * handwritten optimizations
   * stdlib optimizations? `net/http/roundtrip_js.go`, `reflect/*.go`