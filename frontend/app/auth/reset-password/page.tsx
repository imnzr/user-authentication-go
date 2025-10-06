import { ResetPasswordForm } from "@/app/components/reset-password";

export default function ForgotPasswordPage() {
    return (
      <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
        <div className="w-full max-w-sm">
          <ResetPasswordForm />
        </div>
      </div>
    )
  }
  