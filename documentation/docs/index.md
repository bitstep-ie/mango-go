# mango-go {.hidden}

<figure markdown="span">
    ![Logo](./assets/mango-with-text-black.png#only-light)
    ![Logo](./assets/mango-with-text-white.png#only-dark)
    <figcaption>mango-go</figcaption>
</figure>



The below is the list of packages mango4go currently offers. Each one links to its dedicated deep-dive.

- [env](packages/env.md) – environment variable helpers (`StringDefault`, `String`, `MustString`, `IntDefault`, `AsDefault`, etc.)
- [io](packages/io.md) – delete/backup/restore files by extension
- [logger](packages/logger.md) – slog handler with CLI/file/syslog sinks
- [net](packages/net.md) – network helpers
- [random](packages/random.md) – math/crypto random values, password generation, dates
- [slices](packages/slices.md) – generic slice utilities (contains, chunk, unique, …)
- [testutils](packages/testutils.md) – temp files and assertions for tests
- [time](packages/time.md) – start/end-of-day, duration parsing, “time ago” strings

Deprecated env helpers (`EnvOrDefault`, `MustEnv`, `EnvAsInt`, `MustEnvAsInt`, `EnvAsBool`, `MustEnvAsBool`) are scheduled for removal in `v1.0.0`.

You want more help? See the full [guide](guide.md).
Furthermore, you can find on [Discussions forum](https://github.com/bitstep-ie/mango-go/discussions)