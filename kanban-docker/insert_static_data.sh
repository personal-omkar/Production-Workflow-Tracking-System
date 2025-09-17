#!/bin/bash

# Database connection details
DB_NAME="kanban-master"
DB_USER="vsys-kanban-user"
DB_PASSWORD="NewPassword123"
DB_CONTAINER_NAME="irpl-kanban-db"

# Set PGPASSWORD to avoid password prompt during script execution
export PGPASSWORD="$DB_PASSWORD"

# Insert data into tables within the Docker container
docker exec -i "$DB_CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" <<EOF

-- Insert data into user_roles
INSERT INTO public.user_roles (role_name, description, created_by) VALUES
('Admin', 'Administrator role', 'system'),
('Manager', 'Manager role', 'system'),
('Operator', 'Operator role', 'system');

-- Insert data into users (reference role_id from user_roles)
INSERT INTO public.users (username, email, password_hash, role_id, created_by) VALUES
('admin_user', 'admin@example.com', 'hashed_password1', (SELECT id FROM public.user_roles WHERE role_name = 'Admin'), 'system'),
('manager_user', 'manager@example.com', 'hashed_password2', (SELECT id FROM public.user_roles WHERE role_name = 'Manager'), 'system'),
('operator_user', 'operator@example.com', 'hashed_password3', (SELECT id FROM public.user_roles WHERE role_name = 'Operator'), 'system');

-- Insert data into permissions
INSERT INTO public.permissions (permission_name, description, created_by) VALUES
('View Dashboard', 'Permission to view the dashboard', 'system'),
('Manage Users', 'Permission to manage users', 'system'),
('Approve Orders', 'Permission to approve orders', 'system');

-- Insert data into role_permissions (associate roles with permissions)
INSERT INTO public.role_permissions (role_id, permission_id) VALUES
((SELECT id FROM public.user_roles WHERE role_name = 'Admin'), (SELECT id FROM public.permissions WHERE permission_name = 'Manage Users')),
((SELECT id FROM public.user_roles WHERE role_name = 'Manager'), (SELECT id FROM public.permissions WHERE permission_name = 'Approve Orders')),
((SELECT id FROM public.user_roles WHERE role_name = 'Operator'), (SELECT id FROM public.permissions WHERE permission_name = 'View Dashboard'));

-- Insert data into vendors
INSERT INTO public.vendors (vendor_code, vendor_name, contact_info, address, created_by) VALUES
('V001', 'Vendor A', 'contact@vendorA.com', '123 Vendor St', 'system'),
('V002', 'Vendor B', 'contact@vendorB.com', '456 Vendor Ave', 'system');

-- Insert data into compounds
INSERT INTO public.compounds (compound_name, description, created_by) VALUES
('Compound A', 'Description of Compound A', 'system'),
('Compound B', 'Description of Compound B', 'system');

-- Insert data into prod_process
INSERT INTO public.prod_process (name, link, icon, description, status, line_visibility, created_by) VALUES
('Cutting Process', '/cutting', 'icon-cutting.png', 'Cutting operation', 'Active', TRUE, 'system'),
('Welding Process', '/welding', 'icon-welding.png', 'Welding operation', 'Active', FALSE, 'system');

-- Insert data into prod_line
INSERT INTO public.prod_line (name, icon, description, status, created_by) VALUES
('Cutting Line 1', 'icon-line1.png', 'First cutting line', 'Operational', 'system'),
('Welding Line 1', 'icon-line2.png', 'First welding line', 'Operational', 'system');

-- Insert data into prod_process_line (associate prod_process with prod_line)
INSERT INTO public.prod_process_line (prod_process_id, prod_line_id, "order", created_by) VALUES
((SELECT id FROM public.prod_process WHERE name = 'Cutting Process'), (SELECT id FROM public.prod_line WHERE name = 'Cutting Line 1'), 1, 'system'),
((SELECT id FROM public.prod_process WHERE name = 'Welding Process'), (SELECT id FROM public.prod_line WHERE name = 'Welding Line 1'), 2, 'system');

-- Insert data into kb_root
INSERT INTO public.kb_root (running_no, initial_no, created_by) VALUES
(1, 100, 'system'),
(2, 200, 'system');

-- Insert data into kb_data (reference kb_root_id and compound_id)
INSERT INTO public.kb_data (compound_id, mfg_date_time, demand_date_time, exp_date, cell_no, lot_no, location, kb_root_id, created_by) VALUES
((SELECT id FROM public.compounds WHERE compound_name = 'Compound A'), '2024-11-01 08:00:00', '2024-11-10 08:00:00', '2025-11-01 08:00:00', 'C1', 'L1001', 'Warehouse A', (SELECT id FROM public.kb_root WHERE running_no = 1), 'system'),
((SELECT id FROM public.compounds WHERE compound_name = 'Compound B'), '2024-11-01 08:00:00', '2024-11-15 08:00:00', '2025-11-01 08:00:00', 'C2', 'L1002', 'Warehouse B', (SELECT id FROM public.kb_root WHERE running_no = 2), 'system');

-- Insert data into kb_transaction (reference prod_process_line_id and kb_root_id)
INSERT INTO public.kb_transaction (prod_process_id, status, job_id, kb_root_id, prod_process_line_id, created_by) VALUES
((SELECT id FROM public.prod_process WHERE name = 'Cutting Process'), 'In Progress', 101, (SELECT id FROM public.kb_root WHERE running_no = 1), (SELECT id FROM public.prod_process_line WHERE prod_process_id = (SELECT id FROM public.prod_process WHERE name = 'Cutting Process') AND prod_line_id = (SELECT id FROM public.prod_line WHERE name = 'Cutting Line 1')), 'system'),
((SELECT id FROM public.prod_process WHERE name = 'Welding Process'), 'Completed', 102, (SELECT id FROM public.kb_root WHERE running_no = 2), (SELECT id FROM public.prod_process_line WHERE prod_process_id = (SELECT id FROM public.prod_process WHERE name = 'Welding Process') AND prod_line_id = (SELECT id FROM public.prod_line WHERE name = 'Welding Line 1')), 'system');

-- Insert data into kb_extension (reference kb_root_id and vendor_id)
INSERT INTO public.kb_extension (order_id, code, status, kb_root_id, vendor_id, created_by) VALUES
(1, 'EXT1001', 'Active', (SELECT id FROM public.kb_root WHERE running_no = 1), (SELECT id FROM public.vendors WHERE vendor_code = 'V001'), 'system'),
(2, 'EXT1002', 'Active', (SELECT id FROM public.kb_root WHERE running_no = 2), (SELECT id FROM public.vendors WHERE vendor_code = 'V002'), 'system');

-- Insert data into inventory (reference compound_id)
INSERT INTO public.inventory (compound_id, min_quantity, max_quantity, product_type, description, created_by) VALUES
((SELECT id FROM public.compounds WHERE compound_name = 'Compound A'), 10, 100, 'Required Always', 'Main stock for Compound A', 'system'),
((SELECT id FROM public.compounds WHERE compound_name = 'Compound B'), 5, 50, 'Rarely Required', 'Backup stock for Compound B', 'system');

-- Insert data into user_to_vendor (associate users with vendors)
INSERT INTO public.user_to_vendor (user_id, vendor_id, created_on, modified_on) VALUES
((SELECT id FROM public.users WHERE username = 'admin_user'), (SELECT id FROM public.vendors WHERE vendor_code = 'V001'), NOW(), NOW()),
((SELECT id FROM public.users WHERE username = 'manager_user'), (SELECT id FROM public.vendors WHERE vendor_code = 'V002'), NOW(), NOW());

EOF

# Unset the PGPASSWORD variable for security
unset PGPASSWORD
