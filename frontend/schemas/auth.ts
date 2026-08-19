import * as z from "zod";

export const loginFormSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8, "password length is 8 in minimum"),
});

export const signupFormSchema = loginFormSchema.extend({
  username: z.string().min(3, "name length is 3 in minimum"),
});
