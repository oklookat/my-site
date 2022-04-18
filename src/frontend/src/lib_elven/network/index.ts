import FetchDriver from "$lib_elven/network/fetch"
import { Env } from "$lib_elven/tools/paths"

const apiURL = Env.getAPI()
const Fetchd = new FetchDriver(apiURL.toString())
export default Fetchd

/**
 * little bit stupid explanation how things works
 * 
 * We have 3 entities: SSR, API, components.
 * 
 * Components work together with SSR (right now you are reading this on their territory).
 * 
 * Roughly speaking: components it's site page on your browser.
 * 
 * API could be anywhere. And in our case it does the main job: work with auth, database, etc.
 * 
 * As user you don't see SSR and API - it's two servers.
 * 
 * ------- Auth
 * Case: authorize user by username and password.
 * We send {username, password} request from component to API, and API sends token cookie to us.
 * SSR does not take part.
 * 
 * ------- SSR
 * Case: server-side rendering
 * Why we need it? For Google search robot for example.
 * Without SSR (a.k.a SPA) he get's page faster, but without data, and he need to wait, or no wait.
 * In general, there are misunderstandings, bad for SEO, and summary it's bad for our case.
 * 
 * But with SSR he waiting, then he gets response with data. Similar to classic sites logic.
 * Plus we can change every page meta easy, check auth before send component to user, and do other cool things.
 * 
 * How SSR with auth works. For example, we need to get draft articles:
 * 0. We starting at SSR GET method handler
 * 1. Get token cookie from request
 * 2. Convert cookie value to Authorization header
 * 3. On SSR make request to API with this header
 * 4. Get articles response from API, and pass response JSON from SSR to component
 */