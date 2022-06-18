import FetchDriver from '$elven/network/fetch_driver';
import { getApiURL } from '$elven/tools';

const apiURL = getApiURL().toString();
const Fetchd = new FetchDriver(apiURL);
export default Fetchd;
