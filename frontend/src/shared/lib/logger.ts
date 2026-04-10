//femabras/frontend/src/shared/lib/logger.ts
type LogContext = Record<string, unknown>;

export const logger = {
  error: (message: string, error?: unknown, context?: LogContext) => {
    const logEntry = {
      level: "ERROR",
      timestamp: new Date().toISOString(),
      message,
      error: error instanceof Error ? error.message : String(error),
      stack: error instanceof Error ? error.stack : undefined,
      ...context,
    };

    // Outputs structured JSON for server logs (Vercel/Docker/AWS)
    console.error(JSON.stringify(logEntry));
  },

  warn: (message: string, context?: LogContext) => {
    console.warn(
      JSON.stringify({
        level: "WARN",
        timestamp: new Date().toISOString(),
        message,
        ...context,
      }),
    );
  },

  info: (message: string, context?: LogContext) => {
    console.info(
      JSON.stringify({
        level: "INFO",
        timestamp: new Date().toISOString(),
        message,
        ...context,
      }),
    );
  },
};
