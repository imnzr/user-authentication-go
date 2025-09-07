import { LogoutUser } from "@/app/api/auth";
import { useState } from "react";

export function useLogout() {
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    const logout = async () => {
        setLoading(true)
        setError(null)
        try {
            await LogoutUser()
            window.location.href = "/"
        } catch (err) {
            console.error("Logout error: ", err)
            setError("Logout failed")
        } finally {
            setLoading(true)
        }
    }
    return { logout, loading, error }
}