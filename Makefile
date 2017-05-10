GO = env GOOS=linux GOARCH=arm GOARM=5 go

build :
	$(GO) build -v gear-wheel/ev3Wave

deploy : ev3Wave
	scp ev3Wave robot@10.42.0.107:/home/robot/bin

all : build deploy