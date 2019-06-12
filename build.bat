go build -o client.exe
cd Services
cd cinemahall
go build -o cinemahall.exe
cd ..
cd movie
go build -o movie.exe
cd ..
cd reservation
go build -o reservation.exe
cd ..
cd showing
go build -o showing.exe
cd ..
cd user
go build -o user.exe