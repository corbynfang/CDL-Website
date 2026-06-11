import { createClient } from "@supabase/supabase-js";

const supabaseUrl =
  import.meta.env.VITE_SUPABASE_URL ||
  "https://lgnmkdvwaffzgsgmdnsl.supabase.co";
const supabaseAnonKey =
  import.meta.env.VITE_SUPABASE_ANON_KEY ||
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Imxnbm1rZHZ3YWZmemdzZ21kbnNsIiwicm9sZSI6ImFub24iLCJpYXQiOjE3ODA4NTY5MDEsImV4cCI6MjA5NjQzMjkwMX0.lvV6r1qy_NZecs5sQiNoyQxnmJrjWKYAprsAmStxZPk";

export const supabase = createClient(supabaseUrl, supabaseAnonKey);
