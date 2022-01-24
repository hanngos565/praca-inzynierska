# Praca Inżynierska
## pt. `Implementacja i analiza porównawcza algorytmów do detekcji obiektów.`
<!--[![Go Report Card](https://goreportcard.com/badge/github.com/hanngos565/praca-inzynierska)](https://goreportcard.com/badge/github.com/hanngos565/praca-inzynierska)-->

## Wymagania

- [Docker](https://docs.docker.com/compose/install/)
- [Docker-Compose](https://docs.docker.com/compose/install/)

## Sposób instalacji

1. Pobierz plik `docker-compose.yaml`
2. Jeśli chcesz skorzystać z własnego algorytmu - zamień `services.algorithm.image` na nazwę obrazu z algorytmem spełniającym wymagania

## Sposób aktywacji

1. Otwórz terminal
2. Wejdź do katalogu, w którym znajduje się pobrany plik `docker-compose.yaml`
3. Uruchom program poleceniem:

```bash
docker-compse <-p nazwa_projektu> up -d
```

Zatrzymaj program poleceniem:
```bash
docker-compse <-p nazwa_projektu> stop
```
Uruchom ponownie komendą:
```bash
docker-compse <-p nazwa_projektu> start
```
Zatrzymaj program oraz usuń kontenery poleceniem:
```bash
docker-compse <-p nazwa_projektu> down
```

## Wymagania, jakie musi spełniać kontener z algorytmem
Kontener z algorytmem zawiera:
1. Implementację węzłów końcowych
- [`POST /upload_model`](https://github.com/hanngos565/praca-inzynierska/blob/6768b91c11d8ff3cf87842c851aa510d00d4476c/Mask_RCNN/app/app.py#L19)
- [`POST /demo`](https://github.com/hanngos565/praca-inzynierska/blob/6768b91c11d8ff3cf87842c851aa510d00d4476c/Mask_RCNN/app/app.py#L26)
2. Metodę wysyłającą żądanie do serwera, aby zaktualizował wyniki symualcji [`update(id, results)`](https://github.com/hanngos565/praca-inzynierska/blob/6768b91c11d8ff3cf87842c851aa510d00d4476c/Mask_RCNN/app/update.py#L6)
3. Metodę konwertującą base64 na format obrazu przyjmowanego w funkcji symulacji
4. Domyślny model o nazwie `default`

Dodatkowo symulacje muszą być uruchamiane ASYNCHRONICZNIE.

Przykładowa implementacja powyższych wymagań znajduje się w folderze `Mask-RCNN`