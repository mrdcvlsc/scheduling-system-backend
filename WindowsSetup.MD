# Windows 10/11 Setup (CMD)

> All commands below are for the default Command Prompt (CMD.exe).  
> Paste each line and press Enter; comments start with `REM`.

---

## Requirements

1. **Go** v1.24.0 or above  
2. **Git**  
3. **Node.js** v22.11.0 or above  

*(Download/install these before continuing.)*

---

## Local Release Build (with file-system persistence + login)

1. **Open Command Prompt**, then enter the following commands:

```bat
REM 1. create & enter project folder
mkdir scheduling-system
cd scheduling-system

REM 2. clone the backend and frontend repositories
git clone https://github.com/mrdcvlsc/scheduling-system-frontend.git
git clone https://github.com/mrdcvlsc/scheduling-system-backend.git

REM 3. install frontend deps
cd scheduling-system-frontend
npm install

REM 4. prepare backend
cd ..\scheduling-system-backend

REM copy empty data folder (requires /E to include all subfolders, /I to assume destination is a folder)
xcopy ..\scheduling-system-frontend\scheduling-system-temporary-data-empty .\scheduling-system-temporary-data /E /I

REM 5. build backend & frontend
go build -tags netgo -ldflags "-s -w" -o app.exe
cd ..\scheduling-system-frontend
npm run build

REM copy built "dist" into backend
xcopy dist ..\scheduling-system-backend\dist\*.* /E /I

REM 6. set session-only environment variables
set GIN_MODE=debug
set DEV_MODE=local_release
set AUTH=on
set PORT=3000

REM 32-char secrets (change once, then clear MASTER_KEY)
set SESSION_SECRET=xyzxyzxyzxyzxyzxyzxyzxyzxyzxyzxy
set MASTER_KEY=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa

REM 7. run the backend
cd ..\scheduling-system-backend
app.exe
```

2. Copy the displayed IP address and port from the terminal, then paste it into your browser's address bar.
3. Use the `MASTER_KEY` to log in.
4. Create a new department. Remember the Department ID number, Department Code, and Password.
5. Delete the `DEFAULT` department.
6. Stop the `app` process (usually by pressing `Ctrl+C` in the terminal).
7. Enter the command below in the terminal:

```bash
# remove the master key so it cannot be used by users
REM Stop the app (Ctrl+C), then:
set MASTER_KEY=
app.exe
```
---

## Local Development Build (no login)

Open **two** CMD windows.

### Backend window

```bat
cd scheduling-system\scheduling-system-backend
go build -tags netgo -ldflags "-s -w" -o app.exe

REM session-only variables:
set PORT=3000
set GIN_MODE=debug
set DEV_MODE=
set AUTH=off

REM run:
app.exe
```

### Frontend window

```bat
cd scheduling-system\scheduling-system-frontend
npm run dev
```

Browse to the URL shown (e.g. `http://localhost:5173`).

---

## (Optional) Use MongoDB instead of in-memory

Before running `app.exe`, in CMD:

```bat
set USE_DATABASE=MongoDB
set MONGODB_CONNECTION_STRING=your-Go-MongoDB-connection-string
app.exe
```

---

## Deployment to a Web Service

Configure your host (Azure/AWS/etc.) with these environment variables:

```
AUTH=on
GIN_MODE=release
MASTER_KEY=<32-hex-chars>        REM remove after initial setup
SESSION_SECRET=<32-hex-chars>
USE_DATABASE=MongoDB
MONGODB_CONNECTION_STRING=<your-connection-string>
```
