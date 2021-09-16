import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import * as path from "path"

// https://vitejs.dev/config/
const privateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCeGnfqXKZV4fK8
Bf/ib3P4f9NJt1UjFbqbRPqFcUVK222OY48AC7Wclt7e5Z66uYJC5WcXw5TSsWxI
xA1ZiKclpxJP91QgtLXcXd+u/VzFYrS2T2AKHM30JuRI7h4Ia/HWdL4xodyRuq9B
+klaPK38DP454+1OW7gfgAH+0YYWS/T0+W/Dx4Emh+PiJmO43iLTmTCSCzNHL1k/
6VHhG+UjhCisqb3e+csF23o1pwmDgoETYrNRjo7j0mAvXR6z8xEuwHzSyaOZPfAW
SLpwWrNeEnfGKk83WvydgYf7lNhGgP5K2CCqNcuJ5B3KaIUwBhYe6LhY/70wGVfc
S20BVEYlAgMBAAECggEBAJ4X/4L1ZEqKpouHQEkvg3qgbynOTi7IIpHylsPq8Qwa
gOGaAjtcehz77B5c46z5pRy5ga4Mn3tTeOEGUgcOk30eJDycJ56QtS2HHQ5BQYwK
2INYucNC7FrQrQOUP6nSLsBwemLs5L9AOtGxsbmu5ZcdcvV2SixJoMK24Mfqwbb2
aLgDy3klgNYfsUqdJka/xJyOfirSbaNjma32mJezF/Vjh5hHKidUFyVmsi6fcMQW
2qZpuy4wr3Z8YWXGgKhtJbojfOSCI5vMSD0pHDQ/S68UPjChKMOoCLuCF1SXTcJX
WrIkORehJRF9KNiYsn/9I0aYNz6aHaOCBc2QYFQW7yECgYEAy/yOmmggh/xOiZL3
tDZot0HcVfC/P7LmtTJBD3K6qAWYFyLVAOGfOIP8xp6aflwGK6m8g86ZmNAxW2Cs
J7PSAncgMM0i44ElYXJrlqwNBneOjrDMao9CLbqihfdjjID6Iw4pNGt2P0EuWtPM
g5k+8++uWbNDwA5vW5DNUBzsgHkCgYEAxmrWqSIxUs2lvSoIzyzuX6BF64MuAMSe
9/8T0CK+qo244/R9OohU9PeIPP5XIaO4K1/jsKOMejMBiw4bztG6tV1LBf3VoIbJ
9q6j6rtHEeRS9T2jBBKmQiet/Zp1bDhDFmfhgpqoY1CYoElybwFfvrVUcx0oBygo
ZncmcfjvwA0CgYEAgrbYesF29y0K9sWizWyANc6Sbs/S/fKaPEpzvn56dej/yPOd
Iu9xpDmP5YcE4Sv/GEOwdZYSfwYKWs/cAaUNUl8b1aneTYMYfRVhNDAUOKysQWux
0iby9XK/arFYrVSI9aXBacHjCC2wkqteRoYQ/Sf9xaVGKPAXE8n6hvfB5aECgYBD
9a4GijaYpOVCsFYKF9qhmDS6oNit3KqSSKTRokTqB9avBTurQ6zNpsbcFUVkuT3I
ZXSIOMxRwoyhwzhVE1K+8fBTYgi8jc45FsIAERsbjre3VeAy8exONmcKuM/FPFQG
JZiA8lYAhkdwhxeR28eee18sem0OXnvP9LmPSfNcNQKBgQCwix0yOuFkXBpBYg+1
nbZlBRoGoNkZSfHI+ZHsowxA/1gHwnbDK/fRmxcDv5KwYaTwt700z78GpIAc+L8l
TrmDLG3SRLwpkdUurlecY8N2RT1qfi3OEgrOY7DbxNO01geF06fUGoCLXd8ipsQQ
qgIXxLzqpdLgD/swT9WmGp6CnA==
-----END PRIVATE KEY-----
`


const certificate = `
-----BEGIN CERTIFICATE-----
MIIESDCCArCgAwIBAgIRAMkFK44ezbas7POxGT5MIIEwDQYJKoZIhvcNAQELBQAw
ezEeMBwGA1UEChMVbWtjZXJ0IGRldmVsb3BtZW50IENBMSgwJgYDVQQLDB9DT01Q
VVRFUlxzdGVsc0Bjb21wdXRlciAoMDAgMDApMS8wLQYDVQQDDCZta2NlcnQgQ09N
UFVURVJcc3RlbHNAY29tcHV0ZXIgKDAwIDAwKTAeFw0yMTA5MTUxOTUwMDdaFw0y
MzEyMTUxOTUwMDdaMFMxJzAlBgNVBAoTHm1rY2VydCBkZXZlbG9wbWVudCBjZXJ0
aWZpY2F0ZTEoMCYGA1UECwwfQ09NUFVURVJcc3RlbHNAY29tcHV0ZXIgKDAwIDAw
KTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJ4ad+pcplXh8rwF/+Jv
c/h/00m3VSMVuptE+oVxRUrbbY5jjwALtZyW3t7lnrq5gkLlZxfDlNKxbEjEDVmI
pyWnEk/3VCC0tdxd3679XMVitLZPYAoczfQm5EjuHghr8dZ0vjGh3JG6r0H6SVo8
rfwM/jnj7U5buB+AAf7RhhZL9PT5b8PHgSaH4+ImY7jeItOZMJILM0cvWT/pUeEb
5SOEKKypvd75ywXbejWnCYOCgRNis1GOjuPSYC9dHrPzES7AfNLJo5k98BZIunBa
s14Sd8YqTzda/J2Bh/uU2EaA/krYIKo1y4nkHcpohTAGFh7ouFj/vTAZV9xLbQFU
RiUCAwEAAaNvMG0wDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMB
MB8GA1UdIwQYMBaAFH59ZOntVsXl1aOExVaCbN1vYRUxMCUGA1UdEQQeMByCC29r
bG9va2F0LnJ1gg0qLm9rbG9va2F0LnJ1MA0GCSqGSIb3DQEBCwUAA4IBgQCyh9oB
eUu8OBU5SZiGxpgJbH3GYKiG0yvX76bbKKhQ6cELPrVrh2Z6L1MkLpU863TrjPER
EjrY5PqUhdtdco60TaJdIqohsMp43zlPiJwmsB0G+FVktwJwrlXrgaeH6PrkfPqB
RtN0lCRpFqJyxN0kgTiQgI3Xmqjy0Vsy3xV/sH31kTmKCqcWyLw72Wut5s5XsYLT
eYpjrTs2UfesfvXgRYtXQDwdgMP4nNT1aQpWcdemBaFisz0ZiZu4PZFlrFy/O0ks
CL/NcGMCGdnIw/oFnY9xYTL/x1tu93IoG6361+pwd5caxt3/MQp08CIdYrch//pU
0uL7D0DYQfoDu2wjpViSq+Pmo5gaCwSnh2aM/OmW50JzoPHgaB1EYwF2zb4sihCk
NaowhAdkxGdSmb0/qzediTMnlhVo6Zaqd+TtVPyTjf0AFuaawaEwPY3Yw8waeQ/q
Z48nWsmae5XpO7/AhTEXQyA2nldoYK0YkrRRqhRm2a9vbtu0D3HvrvT8fN4=
-----END CERTIFICATE-----
`

const credentials = {key: privateKey, cert: certificate}

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: [
      { find: '@', replacement: path.resolve(__dirname, 'src') },
    ],
  },
  server: {
    open: true,
    https: credentials,
  },
})
