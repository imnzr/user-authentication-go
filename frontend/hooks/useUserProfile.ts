import { GetUserProfile} from "@/app/api/auth";
import { useEffect, useState } from "react";

type User = {
    username: string,
    email: string,
    avatar?: string
}

export function useUserProfile() {
    const [user, setUser] = useState<User | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        async function fetchData() {
            try {
                const data = await GetUserProfile()
                setUser(data)
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            } catch (error: any) {
                setError(error.message)
            } finally {
                setLoading(false)
            }
        }
        fetchData()
    }, [])
    return {user, loading, error}
}