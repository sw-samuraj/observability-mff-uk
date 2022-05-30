# Observability seminar on MFF UK

This repository contains a template and an assignment for the practical part
of the _Observability_ seminar on [MFF UK](https://www.mff.cuni.cz/).

**Following text is intended just for the purpose of the seminar and is
written in the Czech language.**

## Předpoklady

### Nutné

Pro splnění tohoto úkolu je potřeba mít nainstalované následující nástroje:

* Git
* Golang `v1.18`
* `ssh` & `scp` (nebo ekvivalent)
* `curl` (nebo ekvivalent)

### Volitelné

Pokud si budete chtít rozchodit projekt lokálně na svém počítači (doporučeno),
je potřeba mít nainstalované následující backendy:

* [Grafana](//grafana.com/) 8.5.2
* [Loki](//grafana.com/oss/loki/) 2.4.2
* [Prometheus](//prometheus.io/) 2.35.0
* [Jaeger](//www.jaegertracing.io/) 1.33.0

_Všechny pomocné skripty (startování backendů atd.) jsou napsány a odladěny
pro UNIXové prostředí. V případě použití jiného operačního systému (Windows)
je potřeba tyto skripty adaptovat (např. pro PowerShell)._

## Zadání

**Předpoklad:** vytvořené VM podle zadání v [prezentaci od Ládi
Dobiáše](//ulita.ms.mff.cuni.cz/mattermost/ar2122ls/pl/4rjn1xhrgpg5dbb33f47n5ju6y).

1. Fork a `git clone` této repozitory.
2. Úprava zdrojových kódů.
3. Build aplikace
4. (Volitelně) Lokální spuštění aplikace se správnými konfiguračními parametry.
5. `scp` aplikace na studentské VM
6. Spuštění aplikace se správnými konfiguračními parametry na VM.
7. Provolání aplikace pomocí `curl` (nebo ekvivalentního nástroje).
8. Kontrola observability dat _Grafaně_.
9. `git commit` & `git push` vašich změn do vámi forknuté repozitory.

## Odevzdání zadání

Splněné zadání se odevzdá formou: zpráva v aplikaci PARG, v kanále `#nswi150-clouddev`.
Zpráva by měla mít formát:

```
@vit-kotacka @ldobias Odevzdání zadaní.
X-Request-Id: <vaše request ID>
compartment: <váš studentský compartment>
```

## Troubleshooting

Pokud budete mít během vypracování zadání nějaký problém, udělejte následující věci:

1. Pošlete notifikaci @vit-kotacka na PARG kanále `#nswi150-clouddev`
2. Pošlete odkaz na svou forknutou repozitory, která obsahuje vaše nejaktuálnější změny.
3. Pošlete konfigurační parametry, se kterými spouštíte aplikaci.
4. Definujte prostředí, ve kterém pracujete.

## Podrobné zadání

### 1. Fork a `git clone` této repozitory.

1. Na GitHubu kliknout na tlačítko [Fork](//github.com/sw-samuraj/observability-mff-uk/fork)
2. V pracovním adresáři spustit příkaz: `git clone https://github.com/<vaše-repo>/observability-mff-uk.git`

### 2. Úprava zdrojových kódů.

Projít zdrojové kódy (soubory `*.go`) a vyřešit všechny `TODO`. (Po vyřešení můžete `TODO`
řádek smazat.)

Pro porozumnění aplikaci a fungování instrumentace je doporučeno spustiti aplikaci
před úravou `TODO` a postupně během úprav a sledovat výsledné chování.

Také je doporučeno řešit vždy jen jeden z pilířů observability, tedy postupně: logging
metriky a tracing. Jednotlivé oblasti poznáte podle suffixu `TODO`, tj. `TODO Logging:`,
`TODO Metrics:` a `TODO Tracing`.

#### Příklad úpravy `TODO`:

Počáteční stav:

```go
// TODO Logging: Enable json formatting for logging. Uncomment following line.
// logrus.SetFormatter(&logrus.JSONFormatter{})
```

Cílový stav:

```go
logrus.SetFormatter(&logrus.JSONFormatter{})
```

### 3. Build aplikace

Aplikace se zbuilduje jednoduchým:

```shell
cd <naklonovaná repo>
go build
```

Výsledkem je binárka `observability`.

### 4. (Volitelně) Lokální spuštění aplikace se správnými konfiguračními parametry.

#### Lokální spuštění bez backendů a downstream servis

Spustit aplikaci bez externích závislostí lze jednoduchým příkazem:

```shell
./observability
```

#### Lokální spuštění s backendy a downstream servisami

Pokud máte lokálně nainstalované všechny backendy (viz _Předpoklady_ -> _Volitelné_), lze spustit
všechny backendy pomocí příkazu:

```shell
./run-backends.sh
```

Aplikace se pak následně spustí příkazem:

```shell
./run-app.sh
```

**!!! DŮLEŽITÉ !!!** Před spuštěním skriptu `run-backends.sh` si zkontrolujte, že následující
proměnné na začátku souboru odpovídají vašemu prostředí a případně je upravte na korektní
hodnotu:

```shell
LOKI_HOME="${HOME}/dev/loki"
PROMETHEUS_HOME="${HOME}/dev/prometheus"
JAEGER_HOME="${HOME}/dev/jaeger"
GRAFANA_HOME="${HOME}/dev/grafana"
```

Stejně tak si zkontrolujte dále ve skriptu, že názvy jednotlivých binárek odpovídají
vašemu prostředí, či platformě. Např.:

```shell
"${LOKI_HOME}/loki-linux-amd64"
```

Obdobně potom zkontrolujte ve skriptu `kill-backends.sh`, že názvy běžících procesů
jednotlivých backendů odpovídají hodnotám v proměnné `APP`.

#### Provolání aplikace

Aplikace běží lokálně na adresse `localhost:4040`. Lze ji provolat pomocí skriptu
`curl-load.sh`, nebo jednoduchým `curl` příkazem (či alternativou):

```shell
curl -v -H "X-Request-ID: test-42" localhost:4040/
```

#### Očekávaný výstup logů

Provolaná aplikace před změnami by měla vrátit na konzoli následující výstup:

```
INFO[2022-05-18T15:45:39+02:00] starting observability app on: 0.0.0.0:4040   app=my-app func=main
INFO[2022-05-18T15:45:43+02:00] writing response with status: 200             app=my-app func=homeHandler
```

Po vyřešení všech `TODO` by se výstup aplikace měl objevit v souboru `_logs/observability.log`:

```json
{"app":"my-app","func":"main","level":"info","msg":"starting observability app on: 0.0.0.0:4040","time":"2022-05-18T16:23:22+02:00"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"getCorrelationId","level":"warning","msg":"header X-Correlation-ID is empty, no correlation id has been provided","requestId":"load-1-1652883805","time":"2022-05-18T16:23:25+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"tracingMiddleware","level":"debug","msg":"starting tracing...","requestId":"load-1-1652883805","time":"2022-05-18T16:23:25+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"metricsMiddleware","level":"debug","msg":"starting metrics...","requestId":"load-1-1652883805","time":"2022-05-18T16:23:25+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"loggingMiddleware","level":"info","msg":"serving request: GET localhost:4040/","requestId":"load-1-1652883805","time":"2022-05-18T16:23:25+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"loggingMiddleware","level":"debug","msg":"user agent: curl/7.74.0","requestId":"load-1-1652883805","time":"2022-05-18T16:23:25+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"callDownstream","level":"info","msg":"calling downstream service: http://localhost:5050","requestId":"load-1-1652883805","time":"2022-05-18T16:23:25+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"callDownstream","level":"info","msg":"downstream service returned http code: 200","requestId":"load-1-1652883805","time":"2022-05-18T16:23:26+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"callDownstream","level":"info","msg":"downstream service returned request id: e52a834d-ad13-469a-81e1-ecce61a22286","requestId":"load-1-1652883805","time":"2022-05-18T16:23:26+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"callDownstream","level":"debug","msg":"downstream service returned correlation id: a00219e6-11cf-4a94-8190-655dff17ad9f","requestId":"load-1-1652883805","time":"2022-05-18T16:23:26+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"homeHandler","level":"info","msg":"writing response with status: 200","requestId":"load-1-1652883805","time":"2022-05-18T16:23:27+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"metricsMiddleware","level":"debug","msg":"closing metrics...","requestId":"load-1-1652883805","time":"2022-05-18T16:23:27+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
{"app":"my-app","correlationId":"a00219e6-11cf-4a94-8190-655dff17ad9f","func":"tracingMiddleware","level":"debug","msg":"closing tracing...","requestId":"load-1-1652883805","time":"2022-05-18T16:23:27+02:00","traceId":"21beda5b8997fe529daae0f776877d62"}
```

### 5. `scp` aplikace na studentské VM

Vámi vytvořené VM by mělo mít _public IP adresu_. Zkompilovanou binárku aplikace tam
nahrajete příkazem:

```shell
scp observability opc@<vaše IP adresa>:~
```

O pushování logů do logging backendu (Loki) se stará aplikace `promtail` - tu je potřeba
také nahrát na vaše VM. Aplikaci si stáhnete zde:
[promtail-linux-amd64.zip](//github.com/grafana/loki/releases/download/v2.4.2/promtail-linux-amd64.zip).
Spolu s ní je také potřeba na VM nahrát její konfigurační soubor:

```shell
scp <cesta k rozbalenému promtail>/promtail-linux-amd64 opc@<vaše IP adresa>:~
scp _config/promtail-config.yaml opc@<vaše IP adresa>:~
```

Přihlašte se pomocí `ssh` na vaše VM a vytvořte v domácím adresáři (kde by měly být všechny
soubory nahrané přes `scp`) vytvořte adresář `_logs`:

```shell
mkdir _logs
```

Vámi vytvořené prostředí by mělo vypadat takto:

```
.
├── _logs
├── observability
├── promtail-config.yaml
└── promtail-linux-amd64
```

### 6. Spuštění aplikace se správnými konfiguračními parametry na VM.

Prvně je potřeba spustit aplikaci `promtail`:

```shell
./promtail-linux-amd64 -config.file promtail-config.yaml &
```

Aplikace poběží na pozadí (i když `stdout` a `stderr` jsou přesměrovány na konzoli).
Následně pustíte vaši aplikaci:

```shell
./observability -d "http://service1.edu.dobias.info:5050" -t "http://grafana.edu.dobias.info:14268/api/traces"
```

Význam jednotlivých parametrů si můžete vypsat příkazem:

```shell
[opc@sw-samuraj ~]$ ./observability -h
Options:
  -d string
    	Downstream URL. Empty string triggers no call to downstream service.
  -h	Print help.
  -n string
    	Application name. (default "my-app")
  -p string
    	Application port. (default "4040")
  -t string
    	Tracing URL. (default "http://localhost:14268/api/traces")
```

### 7. Provolání aplikace pomocí `curl` (nebo ekvivalentního nástroje).

Protože vaše VM má veřejnou IP adresu, můžete aplikaci provolat z vašeho lokálního
počítače:

```shell
curl -v <vaše IP adresa>:4040/
```

Případně pro poslání většího počtu requestů můžete použít skript `curl-load.sh`,
kde upravíte parametry `MAX_REQ` a `URL`.

### 8. Kontrola observability dat v _Grafaně_.

Vizualizaci vašich observability dat si můžete prohlédnout v Grafaně, která běží
na adrese [http://grafana.edu.dobias.info:3000/](//grafana.edu.dobias.info:3000/)

V levém menu vyberte položku _Explore_ a následně vyberte v rozbalovacím menu
jeden z backendů (Loki, Prometheus, Jaeger). Zobrazete si data z vaší aplikace
pomocí následujících dotazů.

Názvy vašich labelů se můžou drobně měnit, proto si upravte dotazy podle vašich hodnot.

#### Loki

```
{job=~"downstream-1-logs|downstream-2-logs|student-app-logs"} | json | requestId="<vaše X-Request-ID"
```

### Prometheus

```
rate(http_requests_total{job=~"service1|service2|student",path="/"}[1m])
```

### Jaeger

Hledejte podle vašeho `X-Tracing-Id` (zadáte ho do pole formuláře _Trace ID_).

### 9. `git commit` & `git push` vašich změn do vámi forknuté repozitory.

Všechny změny, které jste provedli ve zdrojovém kódu vložte do Gitu a publikujte
do své forknuté GitHub repozitory:

```shell
git add .
git commit -m "commit message"
git push
```
