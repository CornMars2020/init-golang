# !/bin/sh

ROOTPATH=`pwd`

tag=`date +'%Y-%m-%d-%H%M%S'`

mkdir -p ${ROOTPATH}/bin
rm -rf ${ROOTPATH}/bin/$1*

go env | egrep "GOARCH|GOOS|GOHOSTARCH|GOHOSTOS"

projects=(
  "example"
)

function build()
{
  PROJ=$1
  NAME=$2
  cd ${ROOTPATH}/cmd/${PROJ} && go build -o ${NAME} . && mv ${NAME} ${ROOTPATH}/bin

  chmod +x ${ROOTPATH}/bin/${NAME}
}

for proj in "${projects[@]}";
do
  if [ "$1" == "" ]; then
    # build ${proj} ${proj//find/replace};
    echo "building ${proj}"
    # build ${proj} ${proj}-${tag}
    build ${proj} ${proj}
  elif [ "$1" == "${proj}" ]; then
    # build ${proj} ${proj//find/replace};
    echo "building ${proj}"
    # build ${proj} ${proj}-${tag}
    build ${proj} ${proj}
  fi
done
