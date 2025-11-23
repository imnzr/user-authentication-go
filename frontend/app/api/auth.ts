const API_URL = "/api/v1"

export type TokenResponse = {
    access_token: string;
    refresh_token: string;
};

export async function RegisterUser(username:string, email:string, password:string) {
    const res = await fetch(`${API_URL}/auth/signup`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({username, email, password}),
        credentials: "include"
    })
    if (!res.ok) {
        throw new Error("Registration failed")
    }
    return res.json()
}


export async function LoginUser(email:string, password:string): Promise<TokenResponse>{
    const res = await fetch(`${API_URL}/auth/signin`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({email, password})
    })
    const data = await res.json().catch(() => ({}))

    if (!res.ok) {
        const errorMessage = data.error || "Login failed"
        throw new Error(errorMessage)
    }

    localStorage.setItem("access_token", data.access_token)
    localStorage.setItem("refresh_token", data.refresh_token)

    return data as TokenResponse
}


export async function GetUserProfile() {
    const token = localStorage.getItem("access_token")
    if (!token) {
        throw new Error("No access token found")
    }
    const res = await fetch(`${API_URL}/auth/profile`, {
        method: "GET",
        headers: {
            "Authorization": `Bearer ${token}`
        }
    })
    if (!res.ok) {
        if (res.status === 401) {
            localStorage.removeItem("access_token")
            localStorage.removeItem("refresh_token")
            throw new Error("Session expired. Please login again")
        }
        const errorData = await res.json().catch(() => ({}))
        const errorMessage = errorData.error || errorData.Error || "Failed to get user profile"
        throw new Error(errorMessage)
    }
    const data = await res.json()

    return {
        username: data.username,
        email: data.email,
        avatar: data.avatar || "/default-image.png"
    }
}

export async function LogoutUser() {
    const token = localStorage.getItem("access_token")
    if (!token) {
        throw new Error("No token found")
    }
    const res = await fetch(`${API_URL}/auth/logout`, {
        method: "POST",
        headers: {
             "Authorization": `Bearer ${token}`
        },
    })
    if (!res.ok) {
        throw new Error("Logout failed")
    }

    localStorage.removeItem("access_token")
}
