import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { fileURLToPath } from 'url';
import { resolve, dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const projectRootDir = resolve(__dirname);

const privateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDvvnwXjmEdvvPf
B4v4Dpl7iGCVkRA1d7P/2z048rEGi18sXzIV+BYZOLCOY2WjGUPU9j5YXFbRSejR
5r1o3WYO665jo41o0HsVI4c67HkGLGBGqSmYeTY9S3kEvEHicWCXU8M+HBUeug9j
MWcfzvHArgYuTRE1VXo/0Zz42MAgb1K8o7Y/GftNHdP6TZTi3APkz3u3sg+jD64P
vSSqWbuV/ed/AMM5qPSzwpWtIUCZf26md6lw1lu3+PERjijNTx7dmBtwVvmUInBR
tcB3qxwKDzI5N41ctrr2XN3dd+5YL7GwQEG1dJ2JH18decL9wySHiABd2xMqj3hN
gjpzfNCxAgMBAAECggEAQ7zpaOp0rCfo1cpXFOZjHVFCIOBU2RaSSCdU3ANOW1AD
uxDGoQhwjf3DfOe+rt+Tl9dfdC3a5dVhL/Vm7PspO+wBYOu0J7pIHIfmTJLF+yWM
qkpfk1IkLhtl/R2H7Agmx3sGPkVG/mCJungZ8BXTnqJM4sYpp40YNnHuTS1I06C0
mrNF0Lqwc/7YYrMEBtg4RKQdNcYfAVuKXabwUAYvPsu1UAKez1a2PMfyCCR9Thpj
xDFusE/+2NoWytBm+x3VZ8GYddb6GFg0KA0CXzyT9t0VSUjCCnMlHgFfbLS3fk8S
T+sJYzY6OIu6keDjRDc8UqSQHn0fsyESqh7YsyEOAQKBgQD1smo2/LQbEii8b87r
zUp1Jy6QOltz+L17dXEbr1Bf5JwwNWt9ERkkgNUkCBTTlzpuKSzU/gEZ7cReOOiY
gRQoo7orDGH3G5tm6G27Z5/8hZIOLDWQoT1v0jk/ghYHb3F0mVPS9NXdezanfrxs
2pT3fmkTdao5TMvQn7aX2M1s4QKBgQD5zCpPWvHFGI5r1tFaNrtO1D3BFuc/6kTD
ciwZeXlvT5/0iWMT1eTim3FoTDcARBhgTal27NEJwjqpXmb3Xuypsh4x/BUyCeUG
5D1h5SGEXaO5hehvE76lePagKr5gJDJ8T+DXytebZKmXETunm4VuxpwS5dLIP+Jn
qnTfAa+N0QKBgQDMdj7hp4YSCRiRR9JAdMSDnamoKIg2iEbW2LLZwRj30/0OOUZL
+9CaigPVSkxakLz1M5q1XPN5ql4BswYM803dS/L/I/zacMTv7Hh0vhVrRZR36t5g
j23bAmzVpuHCKkGWMIqLQWhAOesPk+4l//h7xVyX7xFLQddSs7Z1wtugwQKBgGWA
ZnjUVYVnEwgSxALKVpYTe/xn2BOKl+JrFCxNxeBLFI4C/XVj5puF2jEtNvQut1D6
fa/LTq2OdNO5flFqCvRxYf+hlOM+Tgfc2mxa9KGagxqxpN7twPvwCRbQB4d9bd5B
fvVwgOetoZO29uyGL1Px9B1hsaKZAuBozRCI2NLhAoGARCB0NAwUw0211HB3LRy8
soVNdL8noW6kDobygWt9EPBmdK+c274ySqASoFnSAWVCeEkDiYauJ1cz7Nu0odPV
NXNwteSQexgvFkKFWG1oA7obSmW7MWgI2icZKwFw6EbvDr1K/QnBc/v5Qi0SsdGE
24way2rGrDr/wUv/LwfEC+g=
-----END PRIVATE KEY-----
`

const certificate = `
-----BEGIN CERTIFICATE-----
MIIEczCCAtugAwIBAgIRAIjq9yJCV6MJC63ulcTb2UAwDQYJKoZIhvcNAQELBQAw
gZcxHjAcBgNVBAoTFW1rY2VydCBkZXZlbG9wbWVudCBDQTE2MDQGA1UECwwtREVT
S1RPUC1HNExKS0owXHN0ZWxzQERFU0tUT1AtRzRMSktKMCAoMDAgMDApMT0wOwYD
VQQDDDRta2NlcnQgREVTS1RPUC1HNExKS0owXHN0ZWxzQERFU0tUT1AtRzRMSktK
MCAoMDAgMDApMB4XDTIxMDkyNTEzNDM1OVoXDTIzMTIyNTEzNDM1OVowYTEnMCUG
A1UEChMebWtjZXJ0IGRldmVsb3BtZW50IGNlcnRpZmljYXRlMTYwNAYDVQQLDC1E
RVNLVE9QLUc0TEpLSjBcc3RlbHNAREVTS1RPUC1HNExKS0owICgwMCAwMCkwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDvvnwXjmEdvvPfB4v4Dpl7iGCV
kRA1d7P/2z048rEGi18sXzIV+BYZOLCOY2WjGUPU9j5YXFbRSejR5r1o3WYO665j
o41o0HsVI4c67HkGLGBGqSmYeTY9S3kEvEHicWCXU8M+HBUeug9jMWcfzvHArgYu
TRE1VXo/0Zz42MAgb1K8o7Y/GftNHdP6TZTi3APkz3u3sg+jD64PvSSqWbuV/ed/
AMM5qPSzwpWtIUCZf26md6lw1lu3+PERjijNTx7dmBtwVvmUInBRtcB3qxwKDzI5
N41ctrr2XN3dd+5YL7GwQEG1dJ2JH18decL9wySHiABd2xMqj3hNgjpzfNCxAgMB
AAGjbzBtMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAfBgNV
HSMEGDAWgBRBi1UVO2Wfm8wkN96v6P5UqrK0/TAlBgNVHREEHjAcggtva2xvb2th
dC5ydYINKi5va2xvb2thdC5ydTANBgkqhkiG9w0BAQsFAAOCAYEAQmAomQlmMi3N
o47Zd7Z42gR4NG+xP5RPkzg61Xjke+M2vLPNOy+Lr/1j+Od+EV2CWeD+CgyTWGsj
HHsLDxRx9eEqi2yy5taq6FPj4lgk3Ujxfsi1bPOcfcH0zDpKopRv/hIEraBTJJl0
kFv8TqCDlzRS7EAqeyZ/Xwgcedlb7MylMA0md62KDneN2umB8fKpY1P6eK4FWW22
J+2KTx/2FRSQjSForNqEAH56f5KP9fqAtT20co9DUD1uUvMB/7NBcbXKM82NCWYy
wlra0pd1uK7HZFDKfMOZKjWwoI5jDooTyf82pLEFsh8bu9nEWQ0LKbC+KN3kpYhs
91i7ZjU/SNqC/hnU8RFZe+JbP8gvIg8Tw5T9TEGi2aBXtf547LcixyjZbO/lyvL4
X5y1eStzorQhWEWQ/Px897L89Qa2mn3NBcvzUzNzzNnaPb6a2Uty9EuNqG6IQHAb
ZrdHLQjktAJAlJMG9TDEtZ/woxB5gg9HzXmvGtCg79jbeKbtcpJf
-----END CERTIFICATE-----
`

const credentials = { key: privateKey, cert: certificate }

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  optimizeDeps: { exclude: ["svelte-router-spa"] },
  resolve: {
    alias: [
      { find: '@', replacement: resolve(projectRootDir, './src') },
    ],
  },
  server: {
    open: true,
    https: credentials,
  },
})
