
namespace StreamDeckService
{
    partial class Form1
    {
        /// <summary>
        /// Erforderliche Designervariable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Verwendete Ressourcen bereinigen.
        /// </summary>
        /// <param name="disposing">True, wenn verwaltete Ressourcen gelöscht werden sollen; andernfalls False.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Vom Windows Form-Designer generierter Code

        /// <summary>
        /// Erforderliche Methode für die Designerunterstützung.
        /// Der Inhalt der Methode darf nicht mit dem Code-Editor geändert werden.
        /// </summary>
        private void InitializeComponent()
        {
            this.components = new System.ComponentModel.Container();
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(Form1));
            this.notifyIcon1 = new System.Windows.Forms.NotifyIcon(this.components);
            this.cmTaskbar = new System.Windows.Forms.ContextMenuStrip(this.components);
            this.cmLabel1 = new System.Windows.Forms.ToolStripMenuItem();
            this.toolStripSeparator1 = new System.Windows.Forms.ToolStripSeparator();
            this.cmQuit = new System.Windows.Forms.ToolStripMenuItem();
            this.textBox1 = new System.Windows.Forms.TextBox();
            this.cmTaskbar.SuspendLayout();
            this.SuspendLayout();
            // 
            // notifyIcon1
            // 
            this.notifyIcon1.ContextMenuStrip = this.cmTaskbar;
            this.notifyIcon1.Icon = ((System.Drawing.Icon)(resources.GetObject("notifyIcon1.Icon")));
            // 
            // cmTaskbar
            // 
            this.cmTaskbar.ImageScalingSize = new System.Drawing.Size(20, 20);
            this.cmTaskbar.Items.AddRange(new System.Windows.Forms.ToolStripItem[] {
            this.cmLabel1,
            this.toolStripSeparator1,
            this.cmQuit});
            this.cmTaskbar.Name = "cmTaskbar";
            this.cmTaskbar.Size = new System.Drawing.Size(209, 58);
            // 
            // cmLabel1
            // 
            this.cmLabel1.Enabled = false;
            this.cmLabel1.Font = new System.Drawing.Font("Segoe UI", 9F, System.Drawing.FontStyle.Bold);
            this.cmLabel1.Name = "cmLabel1";
            this.cmLabel1.Size = new System.Drawing.Size(208, 24);
            this.cmLabel1.Text = "ReCoS Streamdeck";
            // 
            // toolStripSeparator1
            // 
            this.toolStripSeparator1.Name = "toolStripSeparator1";
            this.toolStripSeparator1.Size = new System.Drawing.Size(205, 6);
            // 
            // cmQuit
            // 
            this.cmQuit.Name = "cmQuit";
            this.cmQuit.Size = new System.Drawing.Size(208, 24);
            this.cmQuit.Text = "Quit";
            this.cmQuit.Click += new System.EventHandler(this.cmQuit_Click);
            // 
            // textBox1
            // 
            this.textBox1.Location = new System.Drawing.Point(12, 12);
            this.textBox1.Multiline = true;
            this.textBox1.Name = "textBox1";
            this.textBox1.Size = new System.Drawing.Size(320, 216);
            this.textBox1.TabIndex = 1;
            // 
            // Form1
            // 
            this.ClientSize = new System.Drawing.Size(379, 252);
            this.Controls.Add(this.textBox1);
            this.Icon = ((System.Drawing.Icon)(resources.GetObject("$this.Icon")));
            this.Name = "Form1";
            this.Activated += new System.EventHandler(this.Form1_Activated);
            this.cmTaskbar.ResumeLayout(false);
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.NotifyIcon notifyIcon1;
        private System.Windows.Forms.ContextMenuStrip cmTaskbar;
        private System.Windows.Forms.ToolStripMenuItem cmQuit;
        private System.Windows.Forms.ToolStripMenuItem cmLabel1;
        private System.Windows.Forms.ToolStripSeparator toolStripSeparator1;
        private System.Windows.Forms.TextBox textBox1;
    }
}

