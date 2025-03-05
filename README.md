# Profiles Service (Perfiles y Seguidores)

## Descripción
El servicio de usuarios gestiona los perfiles y las relaciones entre usuarios, como seguidores y seguidos. También permite la búsqueda y optimización del acceso a perfiles mediante cacheo en Redis.

## Características
- Creación y actualización de perfiles de usuario.
- Almacenamiento de datos personales (nombre, foto, bio, etc.).
- Gestión de seguidores y seguidos (follow/unfollow).
- Búsqueda de usuarios por nombre o email.
- Cacheo de perfiles en Redis para optimización.

## Tecnologías Utilizadas
- **Lenguaje**: Go
- **Base de Datos**: PostgreSQL
- **Cache**: Redis
- **Orquestación**: Docker

## Instalación
1. Clona el repositorio:
   ```sh
   git clone <repo-url>
   cd user-service
   ```
2. Configura las variables de entorno en `.env`.
3. Construye y ejecuta el servicio con Docker:
   ```sh
   docker-compose up --build -d
   ```

## Endpoints
| Método | Ruta            | Descripción |
|--------|----------------|-------------|
| POST   | /profile         | Crea un nuevo usuario |
| GET    | /users/:id     | Obtiene información de un usuario |
| PUT    | /users/:id     | Actualiza el perfil de un usuario |
| GET    | /users/search  | Busca usuarios por nombre o email |
| POST   | /follow/:id    | Sigue a un usuario |
| DELETE | /follow/:id    | Deja de seguir a un usuario |

## Seguridad
- Autenticación mediante JWT.
- Límites en la búsqueda de usuarios para prevenir abuso.
- Cacheo en Redis para optimizar la carga de perfiles.

## Contribución
1. Realiza un fork del repositorio.
2. Crea una rama con tu feature (`git checkout -b feature-nueva`).
3. Haz commit de tus cambios (`git commit -m 'Agrega nueva funcionalidad'`).
4. Sube la rama (`git push origin feature-nueva`).
5. Abre un Pull Request.

## Licencia
Este proyecto está bajo la licencia MIT.
