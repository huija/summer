#!/bin/bash
cd container/pipeline || exit 1
go test -v -cover || exit 2

cd ../../dbs || exit 1
go test -v -cover || exit 2

cd ../logs || exit 1
go test -v -cover || exit 2

cd ../srv || exit 1
go test -v -cover || exit 2

cd ../utils || exit 1
go test -v -cover || exit 2

cd .. || exit 1
go test -v -cover || exit 2
