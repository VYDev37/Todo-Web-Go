import axios from 'axios';

const axiosInstance = axios.create({
    baseURL: `http://${window.location.hostname}:${import.meta.env.VITE_API_PORT || 8080}/`,
    headers: {
        'Content-Type': 'application/json',
    }
});

axiosInstance.interceptors.response.use(
    response => response,
    async (error) => {
        if (axios.isAxiosError(error)) {
            const customMessage = error.response?.data?.message || error.message || 'Something went wrong.';
            return Promise.reject(new Error(customMessage));
        }

        return Promise.reject(new Error('Unexpected error occurred.'));
    }
);


export default axiosInstance;