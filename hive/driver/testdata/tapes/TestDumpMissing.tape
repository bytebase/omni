-- pin: default apache/hive:4.0.1 sha256:5194161ef50b80875f937dff04936df047cbe894e9727ab8f0606349345b1bd3 rel/release-4.0.1

== TestDumpMissing default#1
> connect
< error
| -
| gohive connect: Error while executing query: TStatus({StatusCode:ERROR_STATUS InfoMessages:[Server-side error; please check HS2 logs.] SqlState:0x ErrorCode:0x ErrorMessage:0x})
