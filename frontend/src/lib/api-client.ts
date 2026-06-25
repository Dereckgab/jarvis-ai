const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

interface RequestOptions {
  method?: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  headers?: Record<string, string>;
  body?: unknown;
  token?: string;
}

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  details?: string;
}

export class APIClient {
  private static instance: APIClient;

  private constructor() {}

  static getInstance(): APIClient {
    if (!APIClient.instance) {
      APIClient.instance = new APIClient();
    }
    return APIClient.instance;
  }

  private getHeaders(token?: string): Record<string, string> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    return headers;
  }

  async request<T>(
    endpoint: string,
    options: RequestOptions = {}
  ): Promise<ApiResponse<T>> {
    const {
      method = "GET",
      headers = {},
      body,
      token,
    } = options;

    const url = `${API_BASE_URL}${endpoint}`;

    try {
      const response = await fetch(url, {
        method,
        headers: {
          ...this.getHeaders(token),
          ...headers,
        },
        body: body ? JSON.stringify(body) : undefined,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        return {
          success: false,
          error: errorData.message || `HTTP ${response.status}`,
          details: errorData.details,
        };
      }

      const data = await response.json();
      return {
        success: true,
        data,
      };
    } catch (error) {
      return {
        success: false,
        error: "Network error",
        details: error instanceof Error ? error.message : String(error),
      };
    }
  }

  // Authentication endpoints
  async login(email: string, password: string): Promise<ApiResponse<{ access_token: string; refresh_token: string }>> {
    return this.request("/auth/login", {
      method: "POST",
      body: { email, password },
    });
  }

  async register(email: string, password: string, name: string): Promise<ApiResponse<{ access_token: string; refresh_token: string }>> {
    return this.request("/auth/register", {
      method: "POST",
      body: { email, password, name },
    });
  }

  async refreshToken(refreshToken: string): Promise<ApiResponse<{ access_token: string; refresh_token: string }>> {
    return this.request("/auth/refresh", {
      method: "POST",
      body: { refresh_token: refreshToken },
    });
  }

  // AI endpoints
  async chatCompletion(message: string, token: string): Promise<ApiResponse<{ response: string }>> {
    return this.request("/ai/chat", {
      method: "POST",
      body: { prompt: message },
      token,
    });
  }

  // System info endpoints
  async getLatestSystemInfo(token: string): Promise<ApiResponse<any>> {
    return this.request("/system-info/latest", {
      method: "GET",
      token,
    });
  }

  async getSystemInfoHistory(limit: number = 10, offset: number = 0, token: string = ""): Promise<ApiResponse<any[]>> {
    return this.request(`/system-info/history?limit=${limit}&offset=${offset}`, {
      method: "GET",
      token,
    });
  }

  // TTS endpoints
  async generateSpeech(text: string, token: string): Promise<Blob | null> {
    const url = `${API_BASE_URL}/tts/generate`;
    try {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ text }),
      });

      if (!response.ok) {
        console.error("TTS generation failed:", response.statusText);
        return null;
      }

      return await response.blob();
    } catch (error) {
      console.error("TTS request error:", error);
      return null;
    }
  }

  // Memory endpoints
  async saveMemory(content: string, token: string): Promise<ApiResponse<any>> {
    return this.request("/memory", {
      method: "POST",
      body: { content },
      token,
    });
  }

  async searchMemories(query: string, token: string): Promise<ApiResponse<any[]>> {
    return this.request("/memory/search", {
      method: "POST",
      body: { query },
      token,
    });
  }

  async getMemories(token: string): Promise<ApiResponse<any[]>> {
    return this.request("/memory", {
      method: "GET",
      token,
    });
  }

  // Game endpoints
  async searchGames(query: string, token: string): Promise<ApiResponse<any[]>> {
    return this.request(`/games/search?query=${encodeURIComponent(query)}`, {
      method: "GET",
      token,
    });
  }

  async getGameDetails(gameID: string, token: string): Promise<ApiResponse<any>> {
    return this.request(`/games/${gameID}`, {
      method: "GET",
      token,
    });
  }
}

export const apiClient = APIClient.getInstance();
