import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { fileURLToPath } from 'url';
import { resolve, dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const projectRootDir = resolve(__dirname);

const privateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC0opa11cwL6Myt
vX2maghcEnebOkita3WIsEH6N17FGC1PJv6eMTfHPUsfuc9QYC7l42PFWaE1pfiw
GpB1DFBzU4N14ziMDceuDQJbdrdSV6rH6NAp9YpKIV9hCWexL7/bfBkBTrUDUb52
kSDu2papXDqF/RIo2B1i0rNB55NkB7TaiEBKcnz5dXNSrn8wN6zHm1AdtiYskIe5
lmMAJz7lysCooHk9TEjb8Rip02L6VMsiRUxRfEcSYM4rtt4pvKegk4xaIAEbqmyk
UGpM8vxMRa0OBefuMrmYoa7w8kDdYSa5Tpyzp6/BKTi9Cn9pp7AdzN0EmsicfUMk
lh+UivADAgMBAAECggEAJxYLcHs7tNQB/hJDrw/AkVO5p9Aby9K10wb4r9DtY86Q
d1EjxU03GZisuce2EVlOrJXgj+KOHJK2VjzZR0qL2fiXOJobMUsGCyZrfSkaD+pC
XRGGgZ3M969y+QdU9aU0aaTXsE28o2Q7x2It8yE/epQnYbjaBhaFfU3EVh1q6jjX
gi1qUNaEmPCeDUIIcVCKNiiZIUjZcx66LGicRnfTuNmop8Xc3d/PWO9XnS491PjW
FrkyIhdAKCpQcBf65DB5B3+Qn30LMK6OgDEOEkCpbXn99/8NPm0GG7UXKXJv4iF1
AS5IpRW2b95/e+3wDhiBbYCbWnI5Hnzbp6U/exB4cQKBgQDhWISbo1PlwzO4BCY4
/RXfum2sNINBq0YwHuHvbOYdsHXIfAwrmnSGT7zll2m3My4rKHr33/1dYc4LcBHS
p+6u7tUvxP/x/vuSB1500krySf/7NPpiPw8enp6y7iWMhRb4INQD9o2MxOX3IMA4
YSRaDATh7yXAT/sfMZZO1L2HGwKBgQDNNQ79OyIs24D8tmwCeqx4LMX16B20NzNR
+DOnYktvJWKEgU68ci35ThLK9GmUnOoYKHNR+6lLlHmxcxW6Upv2a5TgDn88lz6s
a/WL6/CGYwsNf+uZ8sIc+er1ids2a2/or32K1Oy5xBXGrzzK7ETiPuAQsQWMks31
JgRTg75BOQKBgF67lqvTQq7hm/ltDSCCDNMMDZUKUm+l8BDEutirYCX1C97nfzpJ
aVkZZmDigVe51VvhXaic6md6eWNHjsaPQIEb+FVnrO1v8xRVsrjLuyQWRd22TNhG
iDWoGzE/eluTi69uJZMw8BrQn5h/lS9ebyjHyrmsjqNFPF3Fu9mSB8cnAoGAdOYi
8dM79r3czxnEkez8T0GCTEG8mitQCwQPuOzf+CMd6koXZnbNBdaaEhAGDWkOWlDA
0RVCLhIY+SOMXyYsPsaueq0HgqbORz8BFIKF9Kh7ZWdX+c03E3aziMjNnkS57Z37
MJl3eaPgm9y8DwVlgFU9I4UUmi4PMfZERKcqPNECgYEAlaeRVCcaog86NMi+UChG
swz7zeLHMajUMv0MkCDimjSSm7XXk387Axm0nZxrbTNU9otUnTInZk18/Zen/naY
cbp7b5/CGJpIGmSiMTAM8JA+uxHLgf+niAPJGyWnzGJvVk+lqO0KamRzXfUseoEY
qB7mPpMjigLWvZ07OxvZFYI=
-----END PRIVATE KEY-----
`

const certificate = `
-----BEGIN CERTIFICATE-----
MIIESDCCArCgAwIBAgIRANrazHMMGRZohQPp9o898UcwDQYJKoZIhvcNAQELBQAw
ezEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMSgwJgYDVQQLDB9DT01Q
VVRFUlxreXJseUBjb21wdXRlciAoMDAgMDApMS8wLQYDVQQDDCZta2NlcnQgQ09N
UFVURVJca3lybHlAY29tcHV0ZXIgKDAwIDAwKTAeFw0yMTEwMjIxNjU0MTVaFw0y
NDAxMjIxNjU0MTVaMFMxJzAlBgNVBAoTHm1rY2VydCBkZXZlbG9wbWVudCBjZXJ0
aWZpY2F0ZTEoMCYGA1UECwwfQ09NUFVURVJca3lybHlAY29tcHV0ZXIgKDAwIDAw
KTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALSilrXVzAvozK29faZq
CFwSd5s6SK1rdYiwQfo3XsUYLU8m/p4xN8c9Sx+5z1BgLuXjY8VZoTWl+LAakHUM
UHNTg3XjOIwNx64NAlt2t1JXqsfo0Cn1ikohX2EJZ7Evv9t8GQFOtQNRvnaRIO7a
lqlcOoX9EijYHWLSs0Hnk2QHtNqIQEpyfPl1c1KufzA3rMebUB22JiyQh7mWYwAn
PuXKwKigeT1MSNvxGKnTYvpUyyJFTFF8RxJgziu23im8p6CTjFogARuqbKRQakzy
/ExFrQ4F5+4yuZihrvDyQN1hJrlOnLOnr8EpOL0Kf2mnsB3M3QSayJx9QySWH5SK
8AMCAwEAAaNvMG0wDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMB
MB8GA1UdIwQYMBaAFJTG+ywCC0kyoMm4cme7eZja1K24MCUGA1UdEQQeMByCC29r
bG9va2F0LnJ1gg0qLm9rbG9va2F0LnJ1MA0GCSqGSIb3DQEBCwUAA4IBgQBQU4gc
0iUTVkdWfipTwrqJh566QvUbH6RrkemZCZRRBAkNPhpGtfLxRB6RJ1RLOfsxoNfy
dmJEL00lKmzSPXE/XfPr1oq6BVfGUh0rjiSCwQau9JN8AaEwNplQcjdYzA8L6Bm4
UcMuQftoIa5D5xL+3W9WiyerNFP+jcqAcvyHXcl94fYjF95rQ5X1ppsEj1LIS1Zu
ElewzoglgqS20ilwpBa8gpfq2p+TN62YdYcA46ROj1ZmtadJwnYgO+8VcUBfftea
Rm7xsnXMwaPkATLi9IqIeS87UbcqUOrDbMKdeTSJC8xgrctLQe4I105wRtKkSi99
GqWkghe19OvgA660sEsbJSs33cSW6NVlgY248M6f54Uj0HORYvVTJmqqmyQt352J
rnMu60CcElX093lVxidAe7sNZ5xD+WCFbJ1eAtpBbpLkaeTKI2wK7Db+acVOrfeJ
6t5wrnHnYh/OfTnHGhZBXmp+EccljaOOHBT9QsFjf2Bjjpz1Jlf0SS6Z+x0=
-----END CERTIFICATE-----
`

const credentials = { key: privateKey, cert: certificate }

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  if (command === 'serve') {
    return {
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
      }
    }
  } else {
    // command === 'build'
    return {
      plugins: [svelte()],
      optimizeDeps: { exclude: ["svelte-router-spa"] },
      resolve: {
        alias: [
          { find: '@', replacement: resolve(projectRootDir, './src') },
        ],
      }
    }
  }
})