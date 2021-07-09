BACKEND_DIR="backend"
ADMIN_FRONTEND_DIR="elvenfront"

cd $BACKEND_DIR
npm install
cd ..
cd $ADMIN_FRONTEND_DIR
npm install
cd ..
docker-compose build
docker-compose up