import { useEffect, useState, type ReactElement } from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useAuthStore } from "./auth-store";

type RequireAuthProps = {
  children: ReactElement;
};

type Gate = {
  token: string;
  phase: "checking" | "allowed" | "denied";
};

export function RequireAuth({ children }: RequireAuthProps) {
  const token = useAuthStore((state) => state.token);
  const clearToken = useAuthStore((state) => state.clearToken);
  const location = useLocation();
  const [gate, setGate] = useState<Gate>({ token, phase: "checking" });

  if (gate.token !== token) {
    setGate({ token, phase: "checking" });
  }

  useEffect(() => {
    let active = true;
    const controller = new AbortController();
    const probeToken = token;

    const checkSession = async () => {
      try {
        const headers = new Headers();
        if (probeToken) {
          headers.set("Authorization", `Bearer ${probeToken}`);
        }
        const response = await fetch("/api/v1/system/info", {
          method: "GET",
          headers,
          signal: controller.signal,
        });
        if (!active) {
          return;
        }
        if (response.ok) {
          setGate({ token: probeToken, phase: "allowed" });
          return;
        }
        if (response.status === 401 && probeToken) {
          clearToken();
          return;
        }
        setGate({ token: probeToken, phase: "denied" });
      } catch {
        if (!active) {
          return;
        }
        // A transport failure is not proof the stored token was revoked.
        setGate({ token: probeToken, phase: probeToken ? "allowed" : "denied" });
      }
    };

    void checkSession();

    return () => {
      active = false;
      controller.abort();
    };
  }, [token, clearToken]);

  if (gate.token !== token || gate.phase === "checking") {
    return null;
  }

  if (gate.phase === "allowed") {
    return children;
  }

  const next = `${location.pathname}${location.search}`;
  return <Navigate to={`/login?next=${encodeURIComponent(next)}`} replace />;
}
