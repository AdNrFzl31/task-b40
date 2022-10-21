INSERT INTO tb_project(name, start_date, end_date, description, image)
	VALUES ('Dumbways Web App', '20-09-2022', '24-09-2022', 'ini adalah contoh', 'html.png');
	
UPDATE public.tb_project
	SET name='Javascript', description='Javascript adalah.....'
	WHERE id=3;
	
DELETE FROM public.tb_project
	WHERE id=3;
	
	SELECT id,name,description FROM tb_project