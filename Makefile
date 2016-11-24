cmds=$(patsubst cmds/%.go,%,$(wildcard cmds/*.go))


all: $(cmds)


%: cmds/%.go
	go build $<


make clean:
	rm -f $(cmds)


.PHONY: all
