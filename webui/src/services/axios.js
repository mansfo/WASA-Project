import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

// Interceptor to identify users

instance.interceptors.request.use(
	(config) => {
		const id = localStorage.getItem("userid")

		if(id){
			config.headers['Authorization'] = 'Bearer '+id;
		}
		return config
	},
	(error) =>{
		return Promise.reject(error)
	}
)

export default instance;
