import Navbar from "../components/navbar";
import "bootstrap/dist/css/bootstrap.min.css";
import "./globals.css";

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
    return (
      <html lang="en">
          <body>
              <Navbar />
              {children}
          </body>
      </html>
    );
}
