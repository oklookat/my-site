import FetchDriver from "$lib/network/fetch"
import { Env } from "$lib/tools/paths"

const apiURL = Env.getAPI()


const Fetchd = new FetchDriver(apiURL)
export default Fetchd