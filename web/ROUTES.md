# Frontend Routes

Baseado nos recursos da API (ver `API_ROUTES.md`) e no shell desenhado em `ui.md`. SvelteKit com route groups — `(auth)` e `(app)` são apenas organizacionais e não aparecem na URL.

```
/                                        landing (marketing)

(auth)/login
(auth)/register
(auth)/mfa

(app)/p                                  resolver: redireciona pro projeto do usuário, ou mostra
                                          um estado vazio "crie seu primeiro projeto" se não houver
                                          nenhum — sem dashboard/overview (ver ui.md)
(app)/p/[slug]                           redireciona pra ./links (também sem overview por projeto)
(app)/p/[slug]/links                     GET/POST /projects/{slug}/links
(app)/p/[slug]/links/[code]              GET/PUT/DELETE /projects/{slug}/links/{code}
                                          + GET /projects/{slug}/links/{code}/analytics
(app)/p/[slug]/analytics                 GET /projects/{slug}/analytics
(app)/p/[slug]/webhooks                  GET/POST /projects/{slug}/webhooks/
(app)/p/[slug]/webhooks/[id]/deliveries  GET /projects/{slug}/webhooks/{id}/deliveries

(app)/settings/api-keys                  GET/POST/DELETE /api-keys/
(app)/settings/security                  GET /auth/me, /auth/totp/setup|confirm, DELETE /auth/totp

logout                                   limpa os cookies de auth, chama POST /auth/logout, → /login
```

## Notas

- `/p` (sem slug) é uma rota real, não só um redirect — é o que renderiza quando o usuário não tem nenhum projeto ainda.
- Não existe `/projects` (lista) nem `/dashboard`: o seletor "Active Project ▾" do sidebar (ui.md) é a única forma de trocar de projeto, não há tela de overview separada.
- A URL do frontend usa o prefixo curto `/p/[slug]/...`, desacoplado do path da API (`/projects/{slug}/...`) — decisão deliberada por brevidade.
- `GET /projects/` deve ser buscado uma única vez em `(app)/p/+layout.server.ts` e compartilhado com as rotas `[slug]` abaixo, em vez de refetch por página. O projeto ativo em `[slug]` pode ser derivado dessa mesma lista — a API não tem endpoint de GET para um projeto individual (só list/create/update/delete).
- As streams SSE (`/projects/stream`, `/projects/{slug}/stream`, `/projects/{slug}/links/{code}/stream`) não entram nessa árvore por não serem páginas. Nota para o futuro: a API só autentica via `Authorization: Bearer`, que o `EventSource` do browser não consegue enviar — vai ser necessário um proxy `+server.ts` same-origin lendo o cookie no servidor.
