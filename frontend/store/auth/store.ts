import { User } from "@/lib/query/actions/auth/types";
import secureLocalStorage from "react-secure-storage";
import { create } from "zustand";

interface IAuthStore {
  user: User | null;
  setUser: (user: User) => void;
  logout: () => void;
}

export const useAuthStore = create<IAuthStore>((set) => ({
  user: null,
  setUser: (user) => {
    set({
      user: user,
    });
  },
  logout: () => {
    set({
      user: null,
    });
    secureLocalStorage.clear();
  },
}));
