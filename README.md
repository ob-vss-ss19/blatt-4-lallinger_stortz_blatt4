# Prototypische Implementierung eines Reservierungssystems für ein Kino (Blatt 4 Verteilte Softwaresysteme)

## Was macht dieses Tool?
Es werden hier 2 Binaries zur Verfügung gestellt, zum einen services.exe, hier werden Go Micro Services für die Kinosäle, die Filme, die Reservierungen, Aufführungen und die User gestartet,
zum anderen eine client.exe mit der Anfragen an die Services gestellt werden können, wie zum Beispiel hinzufügen eines Films, buchen einer Aufführung etc.. 

## Getting started

-   Zunächst das Github Repository klonen:

    ```
    git clone https://github.com/ob-vss-ss19/blatt-4-lallinger_stortz_blatt4.git blatt4
    ```

-   Wechseln in das Verzeichnis:

    ```
    cd blatt4
    ```

-   Anschließend bauen der Services mit:

    ```
    go build -o services.exe Services/main.go
    ```

-   Bauen des Clients:

    ```
    go build -o client.exe
    ```

-   Starten der Services:

    ```
    ./Services/services.exe
    ```
    
-   Starten des Clients:

    ```
    client.exe fill
    ```
    
Für eine Liste an Befehlen siehe weiter unten Usage.

## Ausführen mit Docker

-   Images bauen:

    ```
    docker build -f services.dockerfile -t services ./
    docker build -f client.dockerfile -t client ./
    ```

-   ein (Docker)-Netzwerk `testnet` erzeugen:

    ```
    docker network create testnet
    ```

-   Starten der Services (Ports 8092-8096 müssen frei sein) im Netzwerk `testnet`:

    ```
    docker run --rm --net testnet server
    ```

-   Starten des Clients (Port 8091 muss frei sein. Für Optionen siehe Usage weiter unten):

    ```
    docker run --rm --net testnet client fill
    ```

## Usage

-   Über den Client können jedem Service (cinema, movie, reservation, showing und user) Daten hinzugefügt (add), gelöscht (delete) und aufgelistet (get) werden.
Einzige Ausnahme bietet hierbei reservation, hier ist es nicht möglich einfach eine Reservierung hinzuzufügen, diese muss zunächst beantragt (request) werden und anschließend gebucht (book).

    ```
    client.exe SERVICE FUNCTION PARAMS
    SERVICE
     cinema
      FUNCTION
      -add PARAMS: name. Example: client.exe cinema add hall1
      -delete PARAMS: name. Example: client.exe cinema delete hall1
      -get: Example: client.exe cinema get
     movie
      FUNCTION
      -add PARAMS: title. Example: client.exe movie add shrek
      -delete PARAMS: title. Example: client.exe movie delete shrek
      -get: Example: client.exe movie get
     reservation
      FUNCTION
      -request PARAMS: user showingID seats. Example: client.exe reservation request sepp 2 4
       Requests a reservation.
      -book PARAMS: reservationID. Example: client.exe reservation book 1
       Books a reservation.
      -delete PARAMS: reservationID. Example: client.exe reservation delete 1
      -get: Example: client.exe reservation get
     showing
      FUNCTION
      -add PARAMS: movie cinema. Example: client.exe showing add shrek hall1
      -delete PARAMS: showingID. Example: client.exe showing delete 4
      -get: Example: client.exe showing get
     user
      FUNCTION
      -add PARAMS: name. Example: client.exe user add sepp
      -delete PARAMS: name. Example: client.exe user delete sepp
      -get: Example: client.exe user get
     fill
      -Fills services with some data. Example: client.exe fill
      ```
