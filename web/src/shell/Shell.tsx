import { useState } from 'react';
import type { PropsWithChildren } from '../types';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { Footer } from './Footer';

interface ShellProps extends PropsWithChildren {
  showSidebar?: boolean;
  showFooter?: boolean;
}

export const Shell: React.FC<ShellProps> = ({
  children,
  showSidebar = true,
  showFooter = true,
}) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);

  const toggleSidebar = (): void => {
    setIsSidebarOpen(prev => !prev);
  };

  return (
    <div className="shell-container">
      <Header onToggleSidebar={toggleSidebar} />

      <div className="shell-content">
        {showSidebar && (
          <Sidebar isOpen={isSidebarOpen} onClose={() => setIsSidebarOpen(false)} />
        )}

        <main className={`shell-main ${isSidebarOpen && showSidebar ? 'with-sidebar' : ''}`}>
          {children}
        </main>
      </div>

      {showFooter && <Footer />}
    </div>
  );
};
