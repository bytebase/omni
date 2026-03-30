SELECT c.name ?? "unknown" AS displayName, c.nickname ?? c.name ?? "anonymous" AS nick FROM c
